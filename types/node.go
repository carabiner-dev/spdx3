// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"errors"
	"slices"
)

var (
	ErrUnsupportedNodeType = errors.New("unsupported node type")

	// ErrIncompatibleNodeType is returned when a node parses correctly but
	// its class is not valid in the field holding it, such as a non-Element
	// inlined in a collection's element or rootElement property.
	ErrIncompatibleNodeType = errors.New("incompatible node type for field")
)

type ID string

type Profile struct {
	Prefix  string
	Classes map[string]Node
}

// Vocabulary is the set of values an SPDX enumerated property admits.
type Vocabulary[T ~string] []T

// Contains reports whether value is a member of the vocabulary.
func (v Vocabulary[T]) Contains(value T) bool {
	return slices.Contains(v, value)
}

// VocabularyValue is implemented by every string type modelling an SPDX
// vocabulary, so a value can be checked without knowing which vocabulary it
// belongs to. The parser uses it to drop values outside their vocabulary,
// and Validate to report them.
type VocabularyValue interface {
	// IsValid reports whether the value is a member of its vocabulary. The
	// empty string is not a member: an absent property is represented by
	// the field being unset, not by an empty value.
	IsValid() bool
}

func (id ID) GetID() string {
	return string(id)
}

// Dispatcher is an object that takes a node JSON data, and returns a
// node with a concrete type. The job of the dispatcher is to detect the
// node type from the JSON and inform the unmarshaller of the type to
// expect.
type Dispatcher interface {
	UnmarshalNode([]byte) (Node, error)
}

type AddressableById interface {
	GetID() string
}

// Node is an interface requiring the base accesor methods of a node
type Node interface {
	GetSPDXID() string
	GetID() string
	GetType() string
	GetName() string
	GetCreationInfo() Node
}

// NodeRef implements Node but the only method that works is ID
// It also implements all descendant marker interfaces to allow it to be
// used in specialized slices like []AgentDescendant
type NodeRef struct {
	ID   string
	Data []byte
}

func (NodeRef) GetSPDXID() string {
	return ""
}

func (nf NodeRef) GetID() string {
	return nf.ID
}

func (NodeRef) GetType() string {
	return ""
}

func (NodeRef) GetName() string {
	return ""
}

func (NodeRef) GetCreationInfo() Node {
	return nil
}

// Marker methods to implement descendant interfaces
// This is hack and we need to resolve how to plug in
// a universal reference.
func (NodeRef) FromElement()         {}
func (NodeRef) FromAgent()           {}
func (NodeRef) FromArtifact()        {}
func (NodeRef) FromRelationship()    {}
func (NodeRef) FromIntegrityMethod() {}

// The licensing markers, so a reference can stand in for a license anywhere
// a license, a license set member or an exception is expected.
func (NodeRef) FromAnyLicenseInfo()    {}
func (NodeRef) FromExtendableLicense() {}
func (NodeRef) FromLicense()           {}
func (NodeRef) FromLicenseAddition()   {}
