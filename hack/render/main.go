// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

// Command render reads SPDX documents and writes them back out through the
// library, producing the output that the conformance checks validate against
// the SPDX project's own tooling.
//
//	go run ./hack/render -out DIR [document...]
//
// With no documents named, every SPDX 3 document under testdata/corpus is
// rendered, keeping its path relative to that directory.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	spdx3 "github.com/carabiner-dev/spdx3"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

const (
	defaultCorpus = "testdata/corpus"

	// authoredName is the document the library builds itself, which the
	// conformance checks validate alongside the rendered corpus. Rendering
	// existing documents cannot exercise the code paths that turn Go values
	// into SPDX ones, since a value read from a document is written back as
	// it came.
	authoredName = "authored.spdx.json"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	out := flag.String("out", "", "directory to write the rendered documents to (required)")
	in := flag.String("in", defaultCorpus, "directory to read documents from when none are named")
	authored := flag.Bool("authored", false, "also write a document built in Go rather than read from a file")
	flag.Parse()

	if *out == "" {
		flag.Usage()
		return fmt.Errorf("-out is required")
	}

	files := flag.Args()
	baseDir := ""
	if len(files) == 0 {
		var err error
		if files, err = findDocuments(*in); err != nil {
			return err
		}
		baseDir = *in
	}
	if len(files) == 0 {
		return fmt.Errorf("no SPDX 3 documents found in %s", *in)
	}

	if *authored {
		target := filepath.Join(*out, authoredName)
		if err := writeAuthored(target); err != nil {
			return err
		}
		fmt.Println(target)
	}

	for _, file := range files {
		target := filepath.Join(*out, filepath.Base(file))
		if baseDir != "" {
			rel, err := filepath.Rel(baseDir, file)
			if err != nil {
				return fmt.Errorf("locating %s: %w", file, err)
			}
			target = filepath.Join(*out, rel)
		}
		if err := render(file, target); err != nil {
			return err
		}
		fmt.Println(target)
	}
	return nil
}

// render parses a document with the library and writes what the library
// renders back out.
func render(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("reading %s: %w", source, err)
	}

	env, err := spdx3.NewParser().Parse(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("parsing %s: %w", source, err)
	}

	buf := &bytes.Buffer{}
	if err := (&spdx3.Renderer{}).Render(env, buf); err != nil {
		return fmt.Errorf("rendering %s: %w", source, err)
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(target), err)
	}
	if err := os.WriteFile(target, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", target, err)
	}
	return nil
}

// writeAuthored builds a document out of Go values and writes it. The
// timestamps come from time.Now(), so they carry the sub-second precision and
// the local offset that SPDX does not allow, and the document is only valid
// if the library converts them.
func writeAuthored(target string) error {
	now := time.Now()

	alice := core.NewPerson("https://example.com/spdx/alice", "Alice")
	tool := core.NewTool("https://example.com/spdx/tool", "spdx3 render")

	pkg := software.NewPackage("https://example.com/spdx/pkg", "example-lib")
	pkg.PackageVersion = "1.4.2"
	pkg.DownloadLocation = "https://example.com/example-lib-1.4.2.tar.gz"
	pkg.PrimaryPurpose = software.SoftwarePurposeLibrary
	pkg.VerifiedUsing = []core.IntegrityMethodDescendant{
		core.NewHash(core.HashAlgorithmSha256,
			"5d41402abc4b2a76b9719d911017c592f0e8b1b45c0f47b09fb8f0e2e0d9c0aa"),
	}
	// Timestamps carrying the precision and the offset of the machine, so
	// that the conversion to what SPDX allows is exercised.
	pkg.BuiltTime = types.NewDateTime(now.Add(-time.Hour))
	pkg.ReleaseTime = types.NewDateTime(now)

	file := software.NewFile("https://example.com/spdx/file", "./src/main.go")

	doc := core.NewSpdxDocument("https://example.com/spdx/document")
	doc.AddRootElement(pkg)
	doc.ProfileConformance = []core.ProfileIdentifierType{
		core.ProfileIdentifierTypeCore, core.ProfileIdentifierTypeSoftware,
	}

	env := spdx3.NewEnvelope()
	env.Graph.AddNode(alice, tool, doc, pkg, file)
	env.Graph.Relate("https://example.com/spdx/describes", doc, core.RelationshipTypeDescribes, pkg)
	env.Graph.Relate("https://example.com/spdx/contains", pkg, core.RelationshipTypeContains, file)

	creation := core.NewCreationInfo(now, alice)
	creation.CreatedUsing = []types.Node{tool}
	env.Graph.SetCreationInfo(creation)

	buf := &bytes.Buffer{}
	if err := (&spdx3.Renderer{}).Render(env, buf); err != nil {
		return fmt.Errorf("rendering the authored document: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o750); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(target), err)
	}
	if err := os.WriteFile(target, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", target, err)
	}
	return nil
}

func findDocuments(dir string) ([]string, error) {
	var found []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		//nolint:gosec // G122: the tree walked is a directory named on the
		// command line, not attacker-controlled input.
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if isSPDX3Document(data) {
			found = append(found, path)
		}
		return nil
	})
	return found, err
}

func isSPDX3Document(data []byte) bool {
	var probe struct {
		Context json.RawMessage   `json:"@context"`
		Graph   []json.RawMessage `json:"@graph"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}
	return len(probe.Graph) > 0 && strings.Contains(string(probe.Context), "spdx.org/rdf/3.")
}
