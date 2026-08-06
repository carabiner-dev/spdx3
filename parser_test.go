// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/stretchr/testify/require"
)

func textractNodeTypes(t *testing.T, graph *Graph) map[string]int {
	t.Helper()
	ret := map[string]int{}
	for _, n := range *graph {
		if _, ok := ret[n.GetType()]; ok {
			ret[n.GetType()]++
		} else {
			ret[n.GetType()] = 1
		}
	}
	return ret
}

func TestRender(t *testing.T) {
	r := Renderer{}
	n := time.Now()
	c := &core.CreationInfo{
		SpecVersion: "3.0.1",
		CreatedBy: []core.AgentDescendant{
			&core.Person{
				Agent: core.Agent{
					Node: core.Node{
						PreNode: base.PreNode{
							ID: "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
						},
					},
				},
			},
		},
		Created: &n,
		PreNode: base.PreNode{
			ID:   "_:creationinfo",
			Type: "CreationInfo",
		},
	}

	e := &Envelope{
		Context: NewContext(ContextURL301),
		Graph: Graph{
			c,
		},
	}

	r.Render(e, os.Stdout)
	// t.Fail()
}

func TestGraphMarshalJSON(t *testing.T) {
	t.Run("Marshal Graph with nested node references", func(t *testing.T) {
		now := time.Now()

		person := &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						ID:     "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name: "John Doe",
				},
			},
		}

		creationInfo := &core.CreationInfo{
			PreNode: base.PreNode{
				ID:     "_:creationinfo",
				Type:   "CreationInfo",
				SPDXID: "SPDXRef-CreationInfo",
			},
			Name:        "Test Creation",
			SpecVersion: "3.0.1",
			CreatedBy:   []core.AgentDescendant{person},
			Created:     &now,
		}

		graph := Graph{creationInfo, person}

		data, err := graph.MarshalJSON()
		require.NoError(t, err)

		// Unmarshal to verify structure
		var result []map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Len(t, result, 2)

		// Verify CreationInfo node
		ciNode := result[0]
		require.Equal(t, "_:creationinfo", ciNode["@id"])
		require.Equal(t, "CreationInfo", ciNode["type"])
		require.Equal(t, "3.0.1", ciNode["specVersion"])

		// Verify that CreatedBy is serialized as an array of SPDXID references
		createdBy, ok := ciNode["createdBy"].([]interface{})
		require.True(t, ok, "createdBy should be an array")
		require.Len(t, createdBy, 1)
		require.Equal(t, "SPDXRef-Person1", createdBy[0])

		// Verify Person node
		personNode := result[1]
		require.Equal(t, "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1", personNode["@id"])
		require.Equal(t, "SPDXRef-Person1", personNode["spdxId"])
		require.Equal(t, "Person", personNode["type"])
		require.Equal(t, "John Doe", personNode["name"])
	})

	t.Run("Marshal empty Graph", func(t *testing.T) {
		graph := Graph{}
		data, err := graph.MarshalJSON()
		require.NoError(t, err)

		var result []interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)
		require.Empty(t, result)
	})
}

func TestRoundtrip(t *testing.T) {
	t.Run("Roundtrip with testdata", func(t *testing.T) {
		// 1. Parse original JSON
		p := NewParser()
		originalData, err := os.ReadFile("testdata/spdx.json")
		require.NoError(t, err)

		env1, err := p.Parse(bytes.NewReader(originalData))
		require.NoError(t, err)
		require.NotNil(t, env1)

		// 2. Marshal back to JSON
		marshaledData, err := json.MarshalIndent(env1, "", "  ")
		require.NoError(t, err)

		// 3. Parse the marshaled JSON
		env2, err := p.Parse(bytes.NewReader(marshaledData))
		require.NoError(t, err)
		require.NotNil(t, env2)

		// 4. Compare structures
		require.Equal(t, len(env1.Graph), len(env2.Graph))
		require.Equal(t, env1.Context, env2.Context)

		for i, node := range env1.Graph {
			require.Equal(t, node.GetType(), env2.Graph[i].GetType(),
				"Node #%d type mismatch", i)
			require.Equal(t, node.GetID(), env2.Graph[i].GetID(),
				"Node #%d ID mismatch", i)
			require.Equal(t, node.GetSPDXID(), env2.Graph[i].GetSPDXID(),
				"Node #%d SPDXID mismatch", i)
		}

		// 5. Verify CreationInfo nested references are preserved
		for _, n := range env2.Graph {
			if n.GetType() != "CreationInfo" {
				continue
			}

			ci, ok := n.(*core.CreationInfo)
			require.True(t, ok)
			require.Len(t, ci.CreatedBy, 1)
			// After roundtrip, CreatedBy should be a NodeRef with the ID preserved
			require.Equal(t, "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
				ci.CreatedBy[0].GetID())
		}
	})
}

func TestNode(t *testing.T) {
	p := NewParser()
	data, err := os.ReadFile("testdata/spdx.json")
	require.NoError(t, err)
	env, err := p.Parse(bytes.NewReader(data))
	require.NoError(t, err)
	require.NotNil(t, env)
	require.NotNil(t, env.Graph)
	require.Equal(t, "CreationInfo", env.Graph[0].GetType())
	require.IsType(t, &core.CreationInfo{}, env.Graph[0])
	require.IsType(t, &core.Person{}, env.Graph[1])
	require.Len(t, env.Graph, 13)
	person := env.Graph[1].(*core.Person)
	require.Len(t, person.ExternalIdentifier, 1)
	require.EqualValues(t, core.ExternalIdentifier{
		Type:                   "ExternalIdentifier",
		ExternalIdentifierType: "email",
		Identifier:             "suriyawa@tcd.ie",
	}, person.ExternalIdentifier[0])

	require.Equal(t, map[string]int{
		"Bom":                               1,
		"CreationInfo":                      1,
		"SpdxDocument":                      1,
		"Organization":                      1,
		"Person":                            1,
		"Relationship":                      4,
		"dataset_DatasetPackage":            1,
		"simplelicensing_LicenseExpression": 1,
		"software_File":                     2,
	}, textractNodeTypes(t, &env.Graph))

	// Check the values of the CreationInfo node
	for _, n := range env.Graph {
		if n.GetType() != "CreationInfo" {
			continue
		}

		ci, ok := n.(*core.CreationInfo)
		require.True(t, ok)

		require.Equal(t, "3.0.1", ci.SpecVersion)
		require.Equal(t, "_:creationinfo", ci.ID)
		crt, err := time.Parse(time.RFC3339, "2024-05-31T00:00:00Z")
		require.NoError(t, err)
		require.Equal(t, &crt, ci.Created)

		// Created by. This checks that the length and the ID match. When
		// finished, we should compare the pointer address of ci.CreatedBy[0]
		// with the Person node in the SBOM, they should be the same.
		// But this is not yet implemented as we need to dedupe the IDs of
		// all referenced nodes.
		require.Len(t, ci.CreatedBy, 1)
		require.Equal(t, "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1", ci.CreatedBy[0].GetID())
	}
}

func TestUnmarshalNodeWithAlias(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		testFunc func(t *testing.T, data []byte)
	}{
		{
			name:     "CreationInfo full object",
			jsonData: `{"spdxId":"id1","@id":"_:ci1","type":"CreationInfo","name":"test","specVersion":"3.0.1"}`,
			testFunc: func(t *testing.T, data []byte) {
				var ci core.CreationInfo
				err := ci.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:ci1", ci.ID)
				require.Equal(t, "CreationInfo", ci.Type)
				require.Equal(t, "test", ci.Name)
				require.Equal(t, "3.0.1", ci.SpecVersion)
			},
		},
		{
			name:     "CreationInfo string reference",
			jsonData: `"_:ci2"`,
			testFunc: func(t *testing.T, data []byte) {
				var ci core.CreationInfo
				err := ci.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:ci2", ci.ID)
			},
		},
		{
			name:     "Person full object",
			jsonData: `{"spdxId":"id1","@id":"_:p1","type":"Person","name":"John Doe","externalIdentifier":[{"type":"ExternalIdentifier","externalIdentifierType":"email","identifier":"john@example.com"}]}`,
			testFunc: func(t *testing.T, data []byte) {
				var p core.Person
				err := p.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:p1", p.ID)
				require.Equal(t, "Person", p.Type)
				require.Equal(t, "John Doe", p.Name)
				require.Len(t, p.ExternalIdentifier, 1)
				require.Equal(t, "john@example.com", p.ExternalIdentifier[0].Identifier)
			},
		},
		{
			name:     "Person string reference",
			jsonData: `"_:p2"`,
			testFunc: func(t *testing.T, data []byte) {
				var p core.Person
				err := p.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:p2", p.ID)
			},
		},
		{
			name:     "File full object with embedded structs",
			jsonData: `{"spdxId":"id1","@id":"_:f1","type":"software_File","name":"test.go","software_downloadLocation":"https://example.com/test.go"}`,
			testFunc: func(t *testing.T, data []byte) {
				var f software.File
				err := f.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:f1", f.ID)
				require.Equal(t, "software_File", f.Type)
				require.Equal(t, "test.go", f.Name)
			},
		},
		{
			name:     "File string reference",
			jsonData: `"_:f2"`,
			testFunc: func(t *testing.T, data []byte) {
				var f software.File
				err := f.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:f2", f.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, []byte(tt.jsonData))
		})
	}
}

// TestVerifiedUsingHash ensures inline integrity methods are dispatched to
// their concrete types on parse and serialized back inline (they have no
// identifier to be referenced by).
func TestVerifiedUsingHash(t *testing.T) {
	doc := `{"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld", "@graph": [
		{
			"type": "software_File",
			"spdxId": "urn:example:file1",
			"name": "data.bin",
			"verifiedUsing": [
				{"type": "Hash", "algorithm": "sha256", "hashValue": "aabbcc"},
				{"type": "PackageVerificationCode", "algorithm": "sha1", "hashValue": "ddeeff"}
			]
		}
	]}`

	env, err := NewParser().Parse(bytes.NewReader([]byte(doc)))
	require.NoError(t, err)
	require.Len(t, env.Graph, 1)

	file, ok := env.Graph[0].(*software.File)
	require.True(t, ok)
	require.Len(t, file.VerifiedUsing, 2)

	hash, ok := file.VerifiedUsing[0].(*core.Hash)
	require.True(t, ok)
	require.Equal(t, core.HashAlgorithm("sha256"), hash.Algorithm)
	require.Equal(t, "aabbcc", hash.HashValue)

	pvc, ok := file.VerifiedUsing[1].(*core.PackageVerificationCode)
	require.True(t, ok)
	require.Equal(t, "ddeeff", pvc.HashValue)

	// The hashes must survive a round trip as inline objects.
	var buf bytes.Buffer
	require.NoError(t, (&Renderer{}).Render(env, &buf))

	env2, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	file2, ok := env2.Graph[0].(*software.File)
	require.True(t, ok)
	require.Len(t, file2.VerifiedUsing, 2)
	hash2, ok := file2.VerifiedUsing[0].(*core.Hash)
	require.True(t, ok)
	require.Equal(t, "aabbcc", hash2.HashValue)
}

// TestParseSpdxId ensures elements are read from and written to the spec's
// spdxId key, and that non-elements (CreationInfo) don't emit an empty one.
func TestParseSpdxId(t *testing.T) {
	f, err := os.Open("testdata/spdx.json")
	require.NoError(t, err)
	defer f.Close() //nolint:errcheck

	env, err := NewParser().Parse(f)
	require.NoError(t, err)

	ids := map[string]bool{}
	for _, n := range env.Graph {
		ids[n.GetSPDXID()] = true
	}
	require.Contains(t, ids, "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1")
	require.Contains(t, ids, "https://spdx.org/spdxdocs/File2-caf55baf-cd02-406a-b7ec-838842ca869f")
	require.Contains(t, ids, "https://spdx.org/spdxdocs/Relationship/describes-de3f5855-bd5a-403f-9aa6-95dd97599c32")

	var buf bytes.Buffer
	require.NoError(t, (&Renderer{}).Render(env, &buf))
	rendered := buf.String()
	require.Contains(t, rendered, `"spdxId": "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1"`)
	require.NotContains(t, rendered, "spdxID")
	require.NotContains(t, rendered, `"spdxId": ""`)
}
