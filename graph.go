// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"fmt"

	"github.com/carabiner-dev/spdx3/dispatch"
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
		//  TODO(puerco): Dedupe IDs of any fields that are nodes by looking up
		// exising nodes with the same ID
		*g = append(*g, n)
	}

	return nil
}
