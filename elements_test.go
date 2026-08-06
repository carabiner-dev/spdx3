// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

func docWithRootElement(rootElement string) string {
	return `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{
				"spdxId": "https://example.com/doc",
				"type": "SpdxDocument",
				"rootElement": [` + rootElement + `]
			}
		]
	}`
}

// TestInlineElementInCollection guards the element/rootElement properties:
// inlined element objects used to panic the parser, because only NodeRef
// implemented the ElementDescendant marker.
func TestInlineElementInCollection(t *testing.T) {
	doc := docWithRootElement(`{
		"spdxId": "https://example.com/pkg",
		"type": "software_Package",
		"name": "inline-package"
	}`)

	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)

	sd, ok := env.Graph[0].(*core.SpdxDocument)
	require.True(t, ok)
	require.Len(t, sd.RootElement, 1)

	// The inlined element keeps its concrete type instead of decaying.
	pkg, ok := sd.RootElement[0].(*software.Package)
	require.True(t, ok)
	require.Equal(t, "inline-package", pkg.GetName())
	require.Equal(t, "https://example.com/pkg", pkg.GetSPDXID())
}

func TestStringRefInCollection(t *testing.T) {
	doc := docWithRootElement(`"https://example.com/pkg"`)

	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)

	sd, ok := env.Graph[0].(*core.SpdxDocument)
	require.True(t, ok)
	require.Len(t, sd.RootElement, 1)
	require.Equal(t, "https://example.com/pkg", sd.RootElement[0].GetID())
}

// TestNonElementInCollection checks that a node which parses fine but is not
// an Element is rejected with an error rather than panicking the parser.
func TestNonElementInCollection(t *testing.T) {
	doc := docWithRootElement(`{
		"@id": "_:creationinfo",
		"type": "CreationInfo",
		"specVersion": "3.0.1"
	}`)

	_, err := NewParser().Parse(strings.NewReader(doc))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrIncompatibleNodeType)
	require.Contains(t, err.Error(), "CreationInfo")
}
