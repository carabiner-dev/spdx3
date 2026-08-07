// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

func TestNewEnvelope(t *testing.T) {
	env := NewEnvelope()

	// It comes bound to the context of the version this library targets.
	require.Equal(t, ContextURL301, env.Context.String())
	require.Equal(t, specVersion301, env.Context.Version())
	require.Empty(t, env.Graph)

	created := time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	creation := &core.CreationInfo{
		PreNode:     base.PreNode{ID: creationInfoID, Type: core.CreationInfoClass},
		SpecVersion: specVersion301,
		Created:     types.NewDateTime(created),
		CreatedBy:   []core.AgentDescendant{types.NodeRef{ID: "https://example.com/alice"}},
	}
	alice := &core.Person{
		Agent: core.Agent{Node: core.Node{
			PreNode:      base.PreNode{SPDXID: "https://example.com/alice", Type: core.PersonClass},
			Name:         "Alice",
			CreationInfo: creation,
		}},
	}
	pkg := &software.Package{
		SoftwareArtifact: software.SoftwareArtifact{Artifact: core.Artifact{Node: core.Node{
			PreNode:      base.PreNode{SPDXID: "https://example.com/pkg", Type: "software_Package"},
			Name:         "example-lib",
			CreationInfo: creation,
		}}},
		PackageVersion: "1.0.0",
	}

	// Nodes can be added one at a time or several at once, and keep the
	// order they were added in.
	env.Graph.AddNode(creation)
	env.Graph.AddNode(alice, pkg)
	require.Len(t, env.Graph, 3)
	require.Same(t, creation, env.Graph[0])
	require.Same(t, alice, env.Graph[1])
	require.Same(t, pkg, env.Graph[2])

	// What it builds is a document that reads back.
	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	require.Contains(t, buf.String(), `"@context": "`+ContextURL301+`"`)

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Len(t, reparsed.Graph, 3)
	require.Equal(t, "example-lib", reparsed.Graph[2].GetName())
}

// An envelope with nothing in it still renders as a document.
func TestNewEnvelopeRendersEmpty(t *testing.T) {
	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(NewEnvelope(), buf))

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Empty(t, reparsed.Graph)
	require.Equal(t, specVersion301, reparsed.Context.Version())
}
