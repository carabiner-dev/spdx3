// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

// Node is the common ancestor of all node types
type Node struct {
	base.PreNode
	Name               string                `json:"name,omitempty"`
	CreationInfo       *CreationInfo         `json:"creationInfo"`
	Comment            string                `json:"comment,omitempty"`
	Description        string                `json:"description,omitempty"`
	Summary            string                `json:"summary,omitempty"`
	Extension          []ExtensionDescendant `json:"extension,omitempty"`
	ExternalIdentifier []ExternalIdentifier  `json:"externalIdentifier,omitempty"`
	ExternalRef        []ExternalRef         `json:"externalRef,omitempty"`
	VerifiedUsing      []IntegrityMethod     `json:"verifiedUsing,omitempty"`
}

func (bn *Node) GetCreationInfo() types.Node {
	return bn.CreationInfo
}

func (bn *Node) GetName() string {
	return bn.Name
}

type ExtensionDescendant interface {
	types.Node
	FromExtension()
}

// Extension is the abstract base class for all SPDX extensions
type Extension struct {
	Node
}

func (ex *Extension) FromExtension() {}

func (e *Extension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}
