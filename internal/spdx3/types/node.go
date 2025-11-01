package types

import "errors"

var ErrUnsupportedNodeType = errors.New("unsupported node type")

type ID string

type Vocabulary[T ~string] []T

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

type RootedNode struct {
	// RootElement []string `json:"rotElement"`
	RootElement []Node `json:"rootElement"`
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
