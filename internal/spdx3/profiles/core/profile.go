package core

import (
	"time"

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

// Agent represents anything with the potential to act on a system
type Agent struct {
	Node
	Summary            string               `json:"summary,omitempty"`
	Extension          []string             `json:"extension,omitempty"`
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier,omitempty"`
	ExternalRef        []string             `json:"externalRef,omitempty"`
	VerifiedUsing      []string             `json:"verifiedUsing,omitempty"`
}

func (a *Agent) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

// Artifact represents a distinct article or unit within the digital domain
type Artifact struct {
	Node
	BuiltTime      *time.Time `json:"builtTime,omitempty"`
	OriginatedBy   []*Agent   `json:"originatedBy,omitempty"`
	ReleaseTime    *time.Time `json:"releaseTime,omitempty"`
	StandardName   []string   `json:"standardName,omitempty"`
	SuppliedBy     *Agent     `json:"suppliedBy,omitempty"`
	SupportLevel   []string   `json:"supportLevel,omitempty"`
	ValidUntilTime *time.Time `json:"validUntilTime,omitempty"`
}

func (a *Artifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
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

// DictionaryEntry represents a key-value pair mapping
type DictionaryEntry struct {
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// IntegrityMethod provides an independently reproducible mechanism for verification
type IntegrityMethod struct {
	Comment string `json:"comment,omitempty"`
}

// PositiveIntegerRange represents a tuple of two positive integers defining a range
type PositiveIntegerRange struct {
	BeginIntegerRange int `json:"beginIntegerRange"`
	EndIntegerRange   int `json:"endIntegerRange"`
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
