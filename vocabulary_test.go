// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
)

// A relationship whose type is not in the vocabulary, and a package with a
// mix of good and bad purposes, so both the scalar and the list case are
// covered.
const nonconformantDoc = `{
	"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
	"@graph": [
		{
			"spdxId": "https://example.com/rel",
			"type": "Relationship",
			"from": "https://example.com/pkg",
			"relationshipType": "totallyMadeUp",
			"to": ["https://example.com/pkg"]
		},
		{
			"spdxId": "https://example.com/pkg",
			"type": "software_Package",
			"name": "example",
			"software_primaryPurpose": "library",
			"software_additionalPurpose": ["source", "nonsense", "archive"]
		}
	]
}`

func TestParseDropsInvalidVocabularyValues(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(nonconformantDoc))
	require.NoError(t, err)

	rel, ok := env.Graph[0].(*core.Relationship)
	require.True(t, ok)
	require.Empty(t, rel.RelationshipType, "a value outside the vocabulary should be dropped")

	pkg, ok := env.Graph[1].(*software.Package)
	require.True(t, ok)
	// The good entries survive; only the bad one is removed.
	require.Equal(t, software.SoftwarePurposeLibrary, pkg.PrimaryPurpose)
	require.Equal(t, []software.SoftwarePurpose{
		software.SoftwarePurposeSource, software.SoftwarePurposeArchive,
	}, pkg.AdditionalPurpose)

	// Nothing invalid is left, so there is nothing to report or to render.
	require.Empty(t, Validate(env))

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	require.NotContains(t, buf.String(), "totallyMadeUp")
	require.NotContains(t, buf.String(), "nonsense")
}

func TestParseKeepsInvalidVocabularyValuesWhenAsked(t *testing.T) {
	env, err := NewParser(WithInvalidVocabularyValues()).Parse(strings.NewReader(nonconformantDoc))
	require.NoError(t, err)

	rel, ok := env.Graph[0].(*core.Relationship)
	require.True(t, ok)
	require.Equal(t, core.RelationshipType("totallyMadeUp"), rel.RelationshipType)

	pkg, ok := env.Graph[1].(*software.Package)
	require.True(t, ok)
	require.Len(t, pkg.AdditionalPurpose, 3)

	// Kept, so the document round-trips as written...
	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	require.Contains(t, buf.String(), "totallyMadeUp")

	// ...and Validate is what tells you it is wrong.
	findings := Validate(env)
	require.Len(t, findings, 2)

	byProperty := map[string]Finding{}
	for _, f := range findings {
		byProperty[f.Property] = f
	}
	require.Contains(t, byProperty, "relationshipType")
	require.Equal(t, "totallyMadeUp", byProperty["relationshipType"].Value)
	require.Equal(t, "https://example.com/rel", byProperty["relationshipType"].NodeID)
	require.Contains(t, byProperty, "software_additionalPurpose")
	require.Equal(t, "nonsense", byProperty["software_additionalPurpose"].Value)
	require.Contains(t, byProperty["relationshipType"].String(), "not a member")
}

// The option is off again for the next parser, since it is carried by the
// parser rather than left set on the package.
func TestVocabularyOptionDoesNotLeak(t *testing.T) {
	_, err := NewParser(WithInvalidVocabularyValues()).Parse(strings.NewReader(nonconformantDoc))
	require.NoError(t, err)

	env, err := NewParser().Parse(strings.NewReader(nonconformantDoc))
	require.NoError(t, err)
	rel, ok := env.Graph[0].(*core.Relationship)
	require.True(t, ok)
	require.Empty(t, rel.RelationshipType)
}

func TestValidateOnConformantDocument(t *testing.T) {
	env, err := NewParser(WithInvalidVocabularyValues()).Parse(strings.NewReader(`{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{
				"spdxId": "https://example.com/rel",
				"type": "Relationship",
				"from": "https://example.com/a",
				"relationshipType": "contains",
				"to": ["https://example.com/b"]
			}
		]
	}`))
	require.NoError(t, err)
	require.Empty(t, Validate(env))
	require.Empty(t, Validate(nil))
}
