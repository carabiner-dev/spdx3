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
	Name               string                      `json:"name,omitempty"`
	CreationInfo       *CreationInfo               `json:"creationInfo"`
	Comment            string                      `json:"comment,omitempty"`
	Description        string                      `json:"description,omitempty"`
	Summary            string                      `json:"summary,omitempty"`
	Extension          []ExtensionDescendant       `json:"extension,omitempty"`
	ExternalIdentifier []ExternalIdentifier        `json:"externalIdentifier,omitempty"`
	ExternalRef        []ExternalRef               `json:"externalRef,omitempty"`
	VerifiedUsing      []IntegrityMethodDescendant `json:"verifiedUsing,omitempty"`
}

func (bn *Node) GetCreationInfo() types.Node {
	return bn.CreationInfo
}

func (bn *Node) GetName() string {
	return bn.Name
}

// SetCreationInfo records how this element came to be, leaving alone one
// that already says so.
func (bn *Node) SetCreationInfo(creation *CreationInfo) {
	if bn.CreationInfo == nil {
		bn.CreationInfo = creation
	}
}

type ExtensionDescendant interface {
	types.Node
	FromExtension()
}

// Extension is the abstract base class for all SPDX extensions. The model
// defines it as a standalone class, not an Element, and gives it no
// properties of its own, so it carries only the identity fields every node
// serializes.
type Extension struct {
	base.PreNode
}

func (ex *Extension) FromExtension() {}

// GetName satisfies types.Node. An Extension has no name in the model.
func (ex *Extension) GetName() string {
	return ""
}

// GetCreationInfo satisfies types.Node. An Extension is not an Element and
// carries no creation information.
func (ex *Extension) GetCreationInfo() types.Node {
	return nil
}

func (e *Extension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}
