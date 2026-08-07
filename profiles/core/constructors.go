// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
)

// The constructors below build an element with its identifier, its type and
// the properties nothing is useful without. Everything else is set on the
// returned value. Going through them means the type discriminator a document
// is read by is always the right one, which a struct literal has to get
// right by hand.

// NewPerson returns a person identified by spdxID.
func NewPerson(spdxID, name string) *Person {
	return &Person{Agent: Agent{Node: newNode(spdxID, PersonClass, name)}}
}

// NewOrganization returns an organization identified by spdxID.
func NewOrganization(spdxID, name string) *Organization {
	return &Organization{Agent: Agent{Node: newNode(spdxID, OrganizationClass, name)}}
}

// NewSoftwareAgent returns a software agent identified by spdxID.
func NewSoftwareAgent(spdxID, name string) *SoftwareAgent {
	return &SoftwareAgent{Agent: Agent{Node: newNode(spdxID, SoftwareAgentClass, name)}}
}

// NewTool returns a tool identified by spdxID.
func NewTool(spdxID, name string) *Tool {
	return &Tool{Node: newNode(spdxID, ToolClass, name)}
}

// NewSpdxDocument returns a document identified by spdxID. Name what it is
// about with AddRootElement.
func NewSpdxDocument(spdxID string) *SpdxDocument {
	return &SpdxDocument{Bundle: Bundle{ElementCollection: ElementCollection{
		Node: newNode(spdxID, SpdxDocumentClass, ""),
	}}}
}

// NewBom returns a bill of materials identified by spdxID.
func NewBom(spdxID string) *Bom {
	return &Bom{Bundle: Bundle{ElementCollection: ElementCollection{
		Node: newNode(spdxID, BomClass, ""),
	}}}
}

// NewRelationship returns a relationship of relType, from one element to one
// or more others.
func NewRelationship(spdxID string, from types.Node, relType RelationshipType, to ...types.Node) *Relationship {
	return &Relationship{
		Node:             newNode(spdxID, RelationshipClass, ""),
		From:             from,
		RelationshipType: relType,
		To:               to,
	}
}

// NewCreationInfo returns the creation information an element carries,
// identified by the blank node label documents conventionally give it.
func NewCreationInfo(created time.Time, createdBy ...AgentDescendant) *CreationInfo {
	return &CreationInfo{
		PreNode:     base.PreNode{ID: "_:creationinfo", Type: CreationInfoClass},
		SpecVersion: SpecVersion,
		Created:     types.NewDateTime(created),
		CreatedBy:   createdBy,
	}
}

// NewHash returns a hash, which is written inline wherever it is used rather
// than being an element of the graph, and so has no identifier.
func NewHash(algorithm HashAlgorithm, value string) *Hash {
	return &Hash{
		IntegrityMethod: IntegrityMethod{PreNode: base.PreNode{Type: HashClass}},
		Algorithm:       algorithm,
		HashValue:       value,
	}
}

func newNode(spdxID, class, name string) Node {
	return Node{
		PreNode: base.PreNode{SPDXID: spdxID, Type: class},
		Name:    name,
	}
}
