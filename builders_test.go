// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

// A document built with the constructors and the graph methods has the shape
// a document is supposed to have, and reads back as what was built.
func TestBuildDocument(t *testing.T) {
	alice := core.NewPerson("https://example.com/alice", "Alice")
	tool := core.NewTool("https://example.com/tool", "generator")
	pkg := software.NewPackage("https://example.com/pkg", "example-lib")
	pkg.PackageVersion = packageVersion
	file := software.NewFile("https://example.com/file", "./main.go")

	doc := core.NewSpdxDocument("https://example.com/document")
	doc.AddRootElement(pkg)

	env := NewEnvelope()
	env.Graph.AddNode(alice, tool, doc, pkg, file)

	describes := env.Graph.Relate("https://example.com/describes",
		doc, core.RelationshipTypeDescribes, pkg)
	env.Graph.Relate("https://example.com/contains",
		pkg, core.RelationshipTypeContains, file)

	// Relate adds what it builds, and hands it back for further setting.
	require.Len(t, env.Graph, 7)
	require.Same(t, describes, env.Graph[5])
	require.Same(t, doc, describes.From)
	require.Equal(t, core.RelationshipTypeDescribes, describes.RelationshipType)
	require.Same(t, pkg, describes.To[0])

	creation := core.NewCreationInfo(time.Now(), alice)
	creation.CreatedUsing = []types.Node{tool}
	env.Graph.SetCreationInfo(creation)

	// It reaches every element, including the ones Relate built, and joins
	// the graph itself.
	require.Len(t, env.Graph, 8)
	require.Same(t, creation, env.Graph[7])
	for _, node := range env.Graph {
		if node == creation {
			continue
		}
		require.Same(t, creation, node.GetCreationInfo(), "%s has no creation info", node.GetSPDXID())
	}

	// The constructors set the type discriminator a document is read by.
	require.Equal(t, "software_Package", pkg.Type)
	require.Equal(t, "Person", alice.Type)
	require.Equal(t, core.SpecVersion, creation.SpecVersion)

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Len(t, reparsed.Graph, 8)

	var readDoc *core.SpdxDocument
	for _, node := range reparsed.Graph {
		if d, ok := node.(*core.SpdxDocument); ok {
			readDoc = d
		}
	}
	require.NotNil(t, readDoc)
	require.Len(t, readDoc.RootElement, 1)
	require.Equal(t, "example-lib", readDoc.RootElement[0].GetName())
}

// Creation information already stated on an element is left alone.
func TestSetCreationInfoKeepsWhatIsThere(t *testing.T) {
	own := core.NewCreationInfo(time.Now())
	own.ID = "_:own"
	alice := core.NewPerson("https://example.com/alice", "Alice")
	alice.CreationInfo = own

	pkg := software.NewPackage("https://example.com/pkg", "lib")

	env := NewEnvelope()
	env.Graph.AddNode(alice, pkg)

	shared := core.NewCreationInfo(time.Now())
	env.Graph.SetCreationInfo(shared)

	require.Same(t, own, alice.CreationInfo, "an element that says how it was made keeps saying it")
	require.Same(t, shared, pkg.CreationInfo)
}

func TestSetCreationInfoIgnoresNil(t *testing.T) {
	env := NewEnvelope()
	env.Graph.AddNode(core.NewPerson("https://example.com/alice", "Alice"))
	env.Graph.SetCreationInfo(nil)
	require.Len(t, env.Graph, 1)
}
