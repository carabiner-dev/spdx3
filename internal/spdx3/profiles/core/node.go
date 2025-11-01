package core

import (
	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
)

type Node struct {
	base.PreNode
	Name               string               `json:"name,omitempty"`
	CreationInfo       *CreationInfo        `json:"creationInfo"`
	Comment            string               `json:"comment,omitempty"`
	Description        string               `json:"description,omitempty"`
	Summary            string               `json:"summary,omitempty"`
	Extension          []string             `json:"extension,omitempty"`
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier,omitempty"`
	ExternalRef        []ExternalRef        `json:"externalRef,omitempty"`
	VerifiedUsing      []IntegrityMethod    `json:"verifiedUsing,omitempty"`
}

func (bn *Node) GetCreationInfo() types.Node {
	return bn.CreationInfo
}

func (bn *Node) GetName() string {
	return bn.Name
}
