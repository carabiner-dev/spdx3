package internal

import (
	"bytes"
	"os"
	"testing"
	"time"

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
	c := &CreationInfo{
		SpecVersion: "3.0.1",
		CreatedBy: []string{
			"https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
		},
		Created: &n,
		PreNode: PreNode{
			ID:   "_:creationinfo",
			Type: "CreationInfo",
		},
		RootedNode: RootedNode{
			RootElement: []string{},
		},
	}

	e := &Envelope{
		Context: "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		Graph: Graph{
			c,
		},
	}

	r.Render(e, os.Stdout)
	t.Fail()
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
	require.IsType(t, &CreationInfo{}, env.Graph[0])
	require.IsType(t, &Person{}, env.Graph[1])
	require.Len(t, env.Graph, 13)
	person := env.Graph[1].(*Person)
	require.Len(t, person.ExternalIdentifier, 1)
	require.EqualValues(t, ExternalIdentifier{
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
}
