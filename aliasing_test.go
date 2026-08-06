// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

// TestParseNoAliasing guards against the dispatcher unmarshaling into the
// shared class prototypes, which made every node of the same type in a graph
// point to a single instance holding the last node's data.
func TestParseNoAliasing(t *testing.T) {
	parse := func() *Envelope {
		f, err := os.Open("testdata/spdx.json")
		require.NoError(t, err)
		defer f.Close() //nolint:errcheck

		env, err := NewParser().Parse(f)
		require.NoError(t, err)
		require.Len(t, env.Graph, 13)
		return env
	}

	env := parse()

	// Every entry in the graph must be a distinct instance.
	seen := map[types.Node]int{}
	for i, n := range env.Graph {
		if j, dup := seen[n]; dup {
			t.Errorf("graph entries #%d and #%d are the same instance (type %q)", j, i, n.GetType())
		}
		seen[n] = i
	}

	// Nodes of the same type must each keep their own data.
	relTypes := []core.RelationshipType{}
	fileNames := []string{}
	for _, n := range env.Graph {
		switch node := n.(type) {
		case *core.Relationship:
			relTypes = append(relTypes, node.RelationshipType)
		case *software.File:
			fileNames = append(fileNames, node.GetName())
		}
	}
	slices.Sort(relTypes)
	require.Equal(t, []core.RelationshipType{
		"contains", "describes", "hasConcludedLicense", "hasDeclaredLicense",
	}, relTypes)
	slices.Sort(fileNames)
	require.Equal(t, []string{"codebook.csv", "data.csv"}, fileNames)

	// A second parse must build fresh nodes, not hand back (or mutate)
	// the instances of the first graph.
	env2 := parse()
	for i := range env.Graph {
		if env.Graph[i] == env2.Graph[i] {
			t.Errorf("graph entry #%d is shared between two parses (type %q)", i, env.Graph[i].GetType())
		}
	}
}
