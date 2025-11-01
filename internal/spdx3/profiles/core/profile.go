package core

import (
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

type Relationship struct {
	Node
	From              string
	To                []string
	RelationshipTypes string `json:"relationshipType"`
}

func (r *Relationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, r, &r.PreNode)
}

type Person struct {
	Node
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, p, &p.PreNode)
}

type Organization struct {
	Node
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

func (o *Organization) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, o, &o.PreNode)
}

type ExternalIdentifier struct {
	Type                   string   `json:"type"` //  "ExternalIdentifier",
	ExternalIdentifierType string   `json:"externalIdentifierType"`
	IssuingAuthority       string   `json:"issuingAuthority"`
	Identifier             string   `json:"identifier"`
	IdentifierLocator      []string `json:"identifierLocator"`
}

type Bom struct {
	Node
	types.RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

func (b *Bom) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}

type SpdxDocument struct {
	Node
	types.RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

func (sd *SpdxDocument) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sd, &sd.PreNode)
}
