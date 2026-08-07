// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"fmt"

	"github.com/carabiner-dev/spdx3/dispatch"
	"github.com/carabiner-dev/spdx3/marshal"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/types"
)

// Envelope
type Envelope struct {
	Context Context `json:"@context"`
	Graph   Graph   `json:"@graph"`
}

// UnmarshalJSON reads a document in either of the two shapes the
// serialization allows: a @graph holding the elements, or a lone element as
// the document root. A single root element is read into a one-node graph,
// so it is rendered back inside a @graph.
func (e *Envelope) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshaling document: %w", err)
	}

	if context, ok := raw["@context"]; ok {
		if err := e.Context.UnmarshalJSON(context); err != nil {
			return fmt.Errorf("unmarshaling @context: %w", err)
		}
	}

	if graph, ok := raw["@graph"]; ok {
		return e.Graph.UnmarshalJSON(graph)
	}

	single, err := json.Marshal([]json.RawMessage{data})
	if err != nil {
		return fmt.Errorf("reading the root element: %w", err)
	}
	return e.Graph.UnmarshalJSON(single)
}

// Graph is the list of nodes a document carries. The receivers are
// deliberately mixed: json.Unmarshaler has to take a pointer to grow the
// slice, while json.Marshaler takes a value so it is found on a Graph value.
//
//nolint:recvcheck // mixing receivers is required, see above
type Graph []types.Node

// UnmarshalJSON unmarshalls the JSONLD graph into nodes typed to their kinds.
func (g *Graph) UnmarshalJSON(data []byte) error {
	dispatcher := dispatch.New()
	list := []json.RawMessage{}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshaling graph: %w", err)
	}

	for i, prenodeData := range list {
		n, err := dispatcher.UnmarshalNode(prenodeData)
		if err != nil {
			return fmt.Errorf("unmarshaling node #%d: %w", i, err)
		}
		// TODO(puerco): Dedupe IDs of the resulting node
		// TODO(puerco): Dedupe IDs of any fields that are nodes by looking up
		// existing nodes with the same ID
		*g = append(*g, n)
	}

	return nil
}

// referenceableIDs returns the identifiers a nested node can be collapsed
// to a reference string with: those of the graph's own nodes, which carry
// the data a reference resolves to, plus the individuals the spec
// predefines, which documents reference but never serialize. A nested node
// identified by anything else is rendered inline, since a reference to it
// would resolve to nothing.
func (g Graph) referenceableIDs() map[string]struct{} {
	ids := make(map[string]struct{}, len(g))
	for _, n := range g {
		if id := n.GetSPDXID(); id != "" {
			ids[id] = struct{}{}
		}
		if id := n.GetID(); id != "" {
			ids[id] = struct{}{}
		}
	}
	for _, iri := range core.IndividualIRIs() {
		ids[iri] = struct{}{}
	}
	for _, iri := range expandedlicensing.IndividualIRIs() {
		ids[iri] = struct{}{}
	}
	return ids
}

// MarshalJSON marshals the graph to JSON.
// Each top-level node in the graph is serialized fully, with any nested nodes
// within them serialized as SPDXID reference strings.
func (g Graph) MarshalJSON() ([]byte, error) {
	if len(g) == 0 {
		return json.Marshal([]interface{}{})
	}

	marshaler := &marshal.NodeMarshaler{
		ReferenceableIDs: g.referenceableIDs(),
	}
	nodeArray := make([]json.RawMessage, 0, len(g))

	for i, node := range g {
		marshaled, err := marshaler.MarshalNode(node)
		if err != nil {
			return nil, fmt.Errorf("marshaling node #%d (type: %s, id: %s): %w",
				i, node.GetType(), node.GetID(), err)
		}
		nodeArray = append(nodeArray, marshaled)
	}

	return json.Marshal(nodeArray)
}
