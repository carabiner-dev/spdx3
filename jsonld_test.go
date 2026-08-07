// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

func docWithNode(node string) string {
	return `{"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld", "@graph": [` + node + `]}`
}

// The SPDX context aliases type to @type, so a document may name a class
// with either. What we write is always the plain spelling the serialization
// asks for.
func TestAtTypeIsReadAndNormalized(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(docWithNode(`{
		"spdxId": "https://example.com/pkg",
		"@type": "software_Package",
		"name": "lib",
		"software_packageVersion": "1.0"
	}`)))
	require.NoError(t, err)

	pkg, ok := env.Graph[0].(*software.Package)
	require.True(t, ok)
	require.Equal(t, "lib", pkg.GetName())
	require.Equal(t, "1.0", pkg.PackageVersion)

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	require.Contains(t, buf.String(), `"type": "software_Package"`)
	require.NotContains(t, buf.String(), `"@type"`)
}

// A property this library does not know is ignored, so a document written
// against a later version of the specification still reads.
func TestUnknownPropertiesAreIgnored(t *testing.T) {
	// software_artifactSize and intendedUse are new in SPDX 3.1.
	env, err := NewParser().Parse(strings.NewReader(docWithNode(`{
		"spdxId": "https://example.com/pkg",
		"type": "software_Package",
		"name": "lib",
		"software_artifactSize": 4096,
		"intendedUse": "testing"
	}`)))
	require.NoError(t, err)

	pkg, ok := env.Graph[0].(*software.Package)
	require.True(t, ok)
	require.Equal(t, "lib", pkg.GetName())
}

// A node naming nothing but itself is complete, not unreadable.
func TestIdentityOnlyNodeParses(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(docWithNode(`{
		"spdxId": "https://example.com/pkg",
		"type": "software_Package"
	}`)))
	require.NoError(t, err)
	require.Len(t, env.Graph, 1)
}

// But a node where nothing at all binds would parse to an empty element,
// which is worth reporting rather than returning. A document written with
// expanded JSON-LD IRIs instead of the names the context defines looks
// exactly like this.
func TestNodeWithNoKnownPropertiesIsReported(t *testing.T) {
	_, err := NewParser().Parse(strings.NewReader(docWithNode(`{
		"spdxId": "https://example.com/pkg",
		"type": "software_Package",
		"https://spdx.org/rdf/3.0.1/terms/Core/name": "lib",
		"https://spdx.org/rdf/3.0.1/terms/Software/downloadLocation": "https://example.com/l.tgz"
	}`)))
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNoKnownProperties)
	// The message names what could not be read.
	require.Contains(t, err.Error(), "terms/Core/name")
	require.Contains(t, err.Error(), "terms/Software/downloadLocation")
}
