package core

import (
	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
)

type Node struct {
	base.PreNode
	Name         string        `json:"name"`
	CreationInfo *CreationInfo `json:"creationInfo"`
	Comment      string        `json:"comment"`
	Description  string        `json:"description"`
}

func (bn *Node) GetCreationInfo() types.Node {
	return bn.CreationInfo
}

func (bn *Node) GetName() string {
	return bn.Name
}
