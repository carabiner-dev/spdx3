// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

// The person is written after the relationship that names it, so resolving
// cannot happen as nodes are read.
const referencesDoc = `{
	"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
	"@graph": [
		{
			"@id": "_:creationinfo",
			"type": "CreationInfo",
			"specVersion": "3.0.1",
			"created": "2024-05-31T00:00:00Z",
			"createdBy": ["https://example.com/alice"]
		},
		{
			"spdxId": "https://example.com/rel",
			"type": "Relationship",
			"creationInfo": "_:creationinfo",
			"from": "https://example.com/pkg",
			"relationshipType": "contains",
			"to": ["https://example.com/file", "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement"]
		},
		{
			"spdxId": "https://example.com/alice",
			"type": "Person",
			"creationInfo": "_:creationinfo",
			"name": "Alice"
		},
		{
			"spdxId": "https://example.com/pkg",
			"type": "software_Package",
			"creationInfo": "_:creationinfo",
			"name": "example-lib"
		},
		{
			"spdxId": "https://example.com/file",
			"type": "software_File",
			"creationInfo": "_:creationinfo",
			"name": "./main.go"
		}
	]
}`

func TestReferencesResolveToNodes(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(referencesDoc))
	require.NoError(t, err)
	require.Len(t, env.Graph, 5)

	creation, ok := env.Graph[0].(*core.CreationInfo)
	require.True(t, ok)
	rel, ok := env.Graph[1].(*core.Relationship)
	require.True(t, ok)
	alice, ok := env.Graph[2].(*core.Person)
	require.True(t, ok)
	pkg, ok := env.Graph[3].(*software.Package)
	require.True(t, ok)
	file, ok := env.Graph[4].(*software.File)
	require.True(t, ok)

	// A reference to a node written later in the graph still resolves.
	require.Same(t, alice, creation.CreatedBy[0])

	// Single-valued and list-valued node properties both resolve.
	require.Same(t, pkg, rel.From)
	require.Same(t, file, rel.To[0])

	// The blank node every element points at is the one in the graph, so
	// reading a node's creation information needs no lookup.
	require.Same(t, creation, pkg.CreationInfo)
	require.Same(t, creation, alice.CreationInfo)

	// A reference the graph has no node for is left as it was: this one
	// names an individual the specification predefines.
	require.Equal(t, core.NoneElementIRI, rel.To[1].GetID())
	_, stillARef := rel.To[1].(types.NodeRef)
	require.True(t, stillARef)
}

// Resolving must not change what a document says, since a nested node is
// written as a reference exactly when the graph carries it.
func TestResolvingKeepsTheDocumentIntact(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(referencesDoc))
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))

	// Each node's data appears once, in the graph, not inlined into every
	// property that names it.
	require.Equal(t, 1, strings.Count(buf.String(), `"name": "Alice"`))
	require.Contains(t, buf.String(), `"from": "https://example.com/pkg"`)

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Len(t, reparsed.Graph, 5)
}

// A graph whose references form a cycle has to terminate.
func TestResolvingCycles(t *testing.T) {
	doc := `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{"spdxId": "https://example.com/a", "type": "Relationship",
			 "from": "https://example.com/b", "relationshipType": "contains",
			 "to": ["https://example.com/a"]},
			{"spdxId": "https://example.com/b", "type": "Relationship",
			 "from": "https://example.com/a", "relationshipType": "contains",
			 "to": ["https://example.com/b"]}
		]
	}`
	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)

	a, ok := env.Graph[0].(*core.Relationship)
	require.True(t, ok)
	b, ok := env.Graph[1].(*core.Relationship)
	require.True(t, ok)
	require.Same(t, b, a.From)
	require.Same(t, a, b.From)
	require.Same(t, a, a.To[0]) // a node may name itself
}

// A node written inline keeps its own identity, and the references inside it
// are resolved too.
func TestResolvingReachesInsideInlineNodes(t *testing.T) {
	doc := `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{"spdxId": "https://example.com/mit", "type": "expandedlicensing_CustomLicense",
			 "simplelicensing_licenseText": "MIT text"},
			{"spdxId": "https://example.com/and", "type": "expandedlicensing_ConjunctiveLicenseSet",
			 "expandedlicensing_member": [
				"https://example.com/mit",
				{"spdxId": "https://example.com/or", "type": "expandedlicensing_DisjunctiveLicenseSet",
				 "expandedlicensing_member": ["https://example.com/mit", "https://example.com/unknown"]}]}
		]
	}`
	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)

	mit, ok := env.Graph[0].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	and, ok := env.Graph[1].(*expandedlicensing.ConjunctiveLicenseSet)
	require.True(t, ok)

	require.Same(t, mit, and.Member[0])

	// The inlined set is not in the graph, so it stays where it was written,
	// but the references it holds resolve.
	or, ok := and.Member[1].(*expandedlicensing.DisjunctiveLicenseSet)
	require.True(t, ok)
	require.Same(t, mit, or.Member[0])
	require.Equal(t, "https://example.com/unknown", or.Member[1].GetID())
}
