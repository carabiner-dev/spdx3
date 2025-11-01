package spdx3

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/software"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
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
		CreatedBy: []types.Node{
			&core.Person{
				Node: core.Node{
					PreNode: base.PreNode{
						ID: "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
					},
				},
			},
		},
		Created: &n,
		PreNode: base.PreNode{
			ID:   "_:creationinfo",
			Type: "CreationInfo",
		},
		RootedNode: types.RootedNode{
			RootElement: []types.Node{},
		},
	}

	e := &Envelope{
		Context: "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		Graph: Graph{
			c,
		},
	}

	r.Render(e, os.Stdout)
	// t.Fail()
}

func TestNode(t *testing.T) {
	p := Parser{}
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
			jsonData: `{"spdxID":"id1","@id":"_:ci1","type":"CreationInfo","name":"test","specVersion":"3.0.1"}`,
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
				require.Equal(t, "ci2", ci.ID)
			},
		},
		{
			name:     "Person full object",
			jsonData: `{"spdxID":"id1","@id":"_:p1","type":"Person","name":"John Doe","externalIdentifier":[{"type":"ExternalIdentifier","externalIdentifierType":"email","identifier":"john@example.com"}]}`,
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
				require.Equal(t, "p2", p.ID)
			},
		},
		{
			name:     "File full object with embedded structs",
			jsonData: `{"spdxID":"id1","@id":"_:f1","type":"software_File","name":"test.go","software_downloadLocation":"https://example.com/test.go"}`,
			testFunc: func(t *testing.T, data []byte) {
				var f software.File
				err := f.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "_:f1", f.ID)
				require.Equal(t, "software_File", f.Type)
				require.Equal(t, "test.go", f.Name)
				require.Equal(t, "https://example.com/test.go", f.DownloadLocation)
			},
		},
		{
			name:     "File string reference",
			jsonData: `"_:f2"`,
			testFunc: func(t *testing.T, data []byte) {
				var f software.File
				err := f.UnmarshalJSON(data)
				require.NoError(t, err)
				require.Equal(t, "f2", f.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, []byte(tt.jsonData))
		})
	}
}
