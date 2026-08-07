// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/stretchr/testify/require"
)

func TestContextForms(t *testing.T) {
	for _, tc := range []struct {
		name    string
		json    string
		mustErr bool
		urls    []string
		version string
	}{
		{
			name:    "string",
			json:    `"https://spdx.org/rdf/3.0.1/spdx-context.jsonld"`,
			urls:    []string{ContextURL301},
			version: specVersion301,
		},
		{
			name:    "array",
			json:    `["https://spdx.org/rdf/3.0.1/spdx-context.jsonld", {"myns": "https://example.com/terms/"}]`,
			urls:    []string{ContextURL301},
			version: specVersion301,
		},
		{
			name: "object",
			json: `{"myns": "https://example.com/terms/"}`,
		},
		{
			name:    "unversioned url",
			json:    `"https://example.com/context.jsonld"`,
			urls:    []string{"https://example.com/context.jsonld"},
			version: "",
		},
		{name: "number", json: `10`, mustErr: true},
		{name: "null", json: `null`, mustErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := Context{}
			err := json.Unmarshal([]byte(tc.json), &ctx)
			if tc.mustErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.urls, ctx.URLs())
			require.Equal(t, tc.version, ctx.Version())

			// Every legal form must render back exactly as it came in.
			out, err := json.Marshal(ctx)
			require.NoError(t, err)
			compacted := &bytes.Buffer{}
			require.NoError(t, json.Compact(compacted, []byte(tc.json)))
			require.JSONEq(t, compacted.String(), string(out))
		})
	}
}

func TestContextVersion(t *testing.T) {
	for url, version := range map[string]string{
		"https://spdx.org/rdf/3.0.1/spdx-context.jsonld":                    specVersion301,
		"https://spdx.org/rdf/3.1/spdx-context.jsonld":                      "3.1",
		"https://spdx.org/rdf/3/spdx-context.jsonld":                        "3",
		"https://spdx.github.io/spdx-spec/v3.1-RC1/rdf/spdx-context.jsonld": "3.1-RC1",
		"https://spdx.org/rdf/3.0.1/spdx-context.jsonld/other":              "",
		"https://example.com/spdx.org/rdf/3.0.1/spdx-context.jsonld":        "",
	} {
		require.Equal(t, version, NewContext(url).Version(), "url: %s", url)
	}
}

func TestParseContext(t *testing.T) {
	t.Run("testdata document", func(t *testing.T) {
		f, err := os.Open("testdata/spdx.json")
		require.NoError(t, err)
		defer f.Close() //nolint:errcheck

		env, err := NewParser().Parse(f)
		require.NoError(t, err)
		require.Equal(t, specVersion301, env.Context.Version())
		require.Equal(t, ContextURL301, env.Context.String())
	})

	t.Run("single root element document", func(t *testing.T) {
		// The serialization allows a lone element as the document root,
		// with no @graph wrapper.
		doc := `{
			"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
			"spdxId": "https://example.com/spdx/Package1",
			"type": "software_Package",
			"name": "lonely-package",
			"software_packageVersion": "1.0.0"
		}`
		env, err := NewParser().Parse(strings.NewReader(doc))
		require.NoError(t, err)
		require.Equal(t, specVersion301, env.Context.Version())
		require.Len(t, env.Graph, 1)

		pkg, ok := env.Graph[0].(*software.Package)
		require.True(t, ok)
		require.Equal(t, "lonely-package", pkg.GetName())
		require.Equal(t, "1.0.0", pkg.PackageVersion)

		// It is rendered back inside a @graph, which reparses to the same
		// single element.
		buf := &bytes.Buffer{}
		require.NoError(t, (&Renderer{}).Render(env, buf))
		require.Contains(t, buf.String(), `"@graph"`)

		reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
		require.NoError(t, err)
		require.Len(t, reparsed.Graph, 1)
		require.Equal(t, "lonely-package", reparsed.Graph[0].GetName())
	})

	t.Run("array context document", func(t *testing.T) {
		doc := `{
			"@context": ["https://spdx.org/rdf/3.0.1/spdx-context.jsonld", {"myns": "https://example.com/terms/"}],
			"@graph": []
		}`
		env, err := NewParser().Parse(strings.NewReader(doc))
		require.NoError(t, err)
		require.Equal(t, specVersion301, env.Context.Version())

		// The array form must survive rendering verbatim.
		buf := &bytes.Buffer{}
		require.NoError(t, (&Renderer{}).Render(env, buf))
		reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
		require.NoError(t, err)
		require.Equal(t, specVersion301, reparsed.Context.Version())
		require.Equal(t, env.Context.URLs(), reparsed.Context.URLs())
	})
}
