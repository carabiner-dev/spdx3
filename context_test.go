// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

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
			version: "3.0.1",
		},
		{
			name:    "array",
			json:    `["https://spdx.org/rdf/3.0.1/spdx-context.jsonld", {"myns": "https://example.com/terms/"}]`,
			urls:    []string{ContextURL301},
			version: "3.0.1",
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
		"https://spdx.org/rdf/3.0.1/spdx-context.jsonld":                    "3.0.1",
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
		require.Equal(t, "3.0.1", env.Context.Version())
		require.Equal(t, ContextURL301, env.Context.String())
	})

	t.Run("array context document", func(t *testing.T) {
		doc := `{
			"@context": ["https://spdx.org/rdf/3.0.1/spdx-context.jsonld", {"myns": "https://example.com/terms/"}],
			"@graph": []
		}`
		env, err := NewParser().Parse(strings.NewReader(doc))
		require.NoError(t, err)
		require.Equal(t, "3.0.1", env.Context.Version())

		// The array form must survive rendering verbatim.
		buf := &bytes.Buffer{}
		require.NoError(t, (&Renderer{}).Render(env, buf))
		reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
		require.NoError(t, err)
		require.Equal(t, "3.0.1", reparsed.Context.Version())
		require.Equal(t, env.Context.URLs(), reparsed.Context.URLs())
	})
}
