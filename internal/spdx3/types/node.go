package types

type ID string

func (id ID) GetID() string {
	return string(id)
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
	ID string
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
