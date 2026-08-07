// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/types"
)

func TestParseIndividuals(t *testing.T) {
	doc := `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{
				"spdxId": "https://example.com/rel",
				"type": "Relationship",
				"from": "https://example.com/pkg",
				"relationshipType": "hasDeclaredLicense",
				"to": ["https://spdx.org/rdf/3.0.1/terms/Licensing/NoAssertion"]
			},
			{
				"spdxId": "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement",
				"type": "IndividualElement",
				"name": "NONE"
			}
		]
	}`
	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)
	require.Len(t, env.Graph, 2)

	// License references resolve against the predefined IRI.
	rel, ok := env.Graph[0].(*core.Relationship)
	require.True(t, ok)
	require.Len(t, rel.To, 1)
	require.Equal(t, expandedlicensing.NoAssertionLicenseIRI, rel.To[0].GetID())

	// An explicit IndividualElement node dispatches to the concrete type.
	ie, ok := env.Graph[1].(*core.IndividualElement)
	require.True(t, ok)
	require.Equal(t, core.NoneElementIRI, ie.GetSPDXID())
	require.Equal(t, "NONE", ie.Name)
}

func TestRenderIndividuals(t *testing.T) {
	env := &Envelope{
		Context: NewContext(ContextURL301),
		Graph: Graph{
			&core.Relationship{
				Node: core.Node{
					PreNode: base.PreNode{
						SPDXID: "https://example.com/rel",
						Type:   core.RelationshipClass,
					},
				},
				From:             types.NodeRef{ID: "https://example.com/pkg"},
				RelationshipType: "hasDeclaredLicense",
				To:               []types.Node{expandedlicensing.NoAssertionLicense},
			},
		},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))

	// The individual must serialize as its IRI reference, never inline.
	require.Contains(t, buf.String(), `"`+expandedlicensing.NoAssertionLicenseIRI+`"`)
	require.NotContains(t, buf.String(), "NOASSERTION")
}
