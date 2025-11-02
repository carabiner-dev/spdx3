// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"fmt"

	"github.com/carabiner-dev/spdx3/dispatch"
	"github.com/carabiner-dev/spdx3/marshal"
	"github.com/carabiner-dev/spdx3/types"
)

// Envelope
type Envelope struct {
	Context string `json:"@context"`
	Graph   Graph  `json:"@graph"`
}

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
		// exising nodes with the same ID
		*g = append(*g, n)
	}

	return nil
}

// MarshalJSON marshals the graph to JSON.
// Each top-level node in the graph is serialized fully, with any nested nodes
// within them serialized as SPDXID reference strings.
func (g Graph) MarshalJSON() ([]byte, error) {
	if len(g) == 0 {
		return json.Marshal([]interface{}{})
	}

	marshaler := &marshal.NodeMarshaler{}
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
