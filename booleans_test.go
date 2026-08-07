// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/stretchr/testify/require"
)

// TestOptionalBooleansRoundTrip guards the tri-state optional booleans: an
// explicit false must survive the round trip instead of being swallowed by
// omitempty, and an absent property must stay absent rather than rendering
// as false.
func TestOptionalBooleansRoundTrip(t *testing.T) {
	doc := `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{
				"spdxId": "https://example.com/license-explicit",
				"type": "expandedlicensing_CustomLicense",
				"simplelicensing_licenseText": "Some license text",
				"expandedlicensing_isFsfLibre": false,
				"expandedlicensing_isOsiApproved": true
			},
			{
				"spdxId": "https://example.com/license-absent",
				"type": "expandedlicensing_CustomLicense",
				"simplelicensing_licenseText": "Other license text"
			}
		]
	}`

	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)
	require.Len(t, env.Graph, 2)

	explicit, ok := env.Graph[0].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.NotNil(t, explicit.IsFsfLibre)
	require.False(t, *explicit.IsFsfLibre)
	require.NotNil(t, explicit.IsOsiApproved)
	require.True(t, *explicit.IsOsiApproved)
	require.Nil(t, explicit.IsDeprecatedLicenseId)

	absent, ok := env.Graph[1].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.Nil(t, absent.IsFsfLibre)
	require.Nil(t, absent.IsOsiApproved)

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	out := buf.String()
	require.Contains(t, out, `"expandedlicensing_isFsfLibre": false`)
	require.Contains(t, out, `"expandedlicensing_isOsiApproved": true`)
	// Unset booleans must not materialize in the output.
	require.NotContains(t, out, "isDeprecatedLicenseId")
	require.Equal(t, 1, strings.Count(out, "isFsfLibre"))

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	reExplicit, ok := reparsed.Graph[0].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.NotNil(t, reExplicit.IsFsfLibre)
	require.False(t, *reExplicit.IsFsfLibre)
	reAbsent, ok := reparsed.Graph[1].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.Nil(t, reAbsent.IsFsfLibre)
}
