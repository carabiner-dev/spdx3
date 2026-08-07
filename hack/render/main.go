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

	spdx3 "github.com/carabiner-dev/spdx3"
)

const defaultCorpus = "testdata/corpus"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	out := flag.String("out", "", "directory to write the rendered documents to (required)")
	in := flag.String("in", defaultCorpus, "directory to read documents from when none are named")
	flag.Parse()

	if *out == "" {
		flag.Usage()
		return fmt.Errorf("-out is required")
	}

	files := flag.Args()
	base := ""
	if len(files) == 0 {
		var err error
		if files, err = findDocuments(*in); err != nil {
			return err
		}
		base = *in
	}
	if len(files) == 0 {
		return fmt.Errorf("no SPDX 3 documents found in %s", *in)
	}

	for _, file := range files {
		target := filepath.Join(*out, filepath.Base(file))
		if base != "" {
			rel, err := filepath.Rel(base, file)
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

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(target), err)
	}
	if err := os.WriteFile(target, buf.Bytes(), 0o644); err != nil {
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
