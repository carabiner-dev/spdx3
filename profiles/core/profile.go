package core

import (
	"time"

	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const Prefix = "core"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"Element":                     &Node{},
		"CreationInfo":                &CreationInfo{},
		"Relationship":                &Relationship{},
		"Agent":                       &Agent{},
		"Artifact":                    &Artifact{},
		"Person":                      &Person{},
		"Organization":                &Organization{},
		"Annotation":                  &Annotation{},
		"SoftwareAgent":               &SoftwareAgent{},
		"Tool":                        &Tool{},
		"LifecycleScopedRelationship": &LifecycleScopedRelationship{},
		"ElementCollection":           &ElementCollection{},
		"Bundle":                      &Bundle{},
		"Bom":                         &Bom{},
		"SpdxDocument":                &SpdxDocument{},
	},
}

type RelationshipDescendant interface {
	FromRelationship()
}

type Relationship struct {
	Node
	From             types.Node               `json:"from"`
	To               []types.Node             `json:"to"`
	RelationshipType RelationshipType         `json:"relationshipType"`
	Completeness     RelationshipCompleteness `json:"completeness,omitempty"`
	StartTime        *time.Time               `json:"startTime,omitempty"`
	EndTime          *time.Time               `json:"endTime,omitempty"`
}

func (r *Relationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, r, &r.PreNode)
}

func (r *Relationship) FromRelationship() {}

// Element is the base element, the node of the SPDX graph
type Element struct {
	Node
}

func (e *Element) FromElement() {}

type ElementDescendant interface {
	types.Node
	FromElement()
}

var _ types.Node = (AgentDescendant)(nil)

type AgentDescendant interface {
	types.Node
	FromAgent()
}

// Agent represents anything with the potential to act on a system
type Agent struct {
	Node
}

func (a *Agent) FromAgent() {}

func (a *Agent) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

// Artifact represents a distinct article or unit within the digital domain
type Artifact struct {
	Node
	BuiltTime      *time.Time    `json:"builtTime,omitempty"`
	OriginatedBy   []*Agent      `json:"originatedBy,omitempty"`
	ReleaseTime    *time.Time    `json:"releaseTime,omitempty"`
	StandardName   []string      `json:"standardName,omitempty"`
	SuppliedBy     *Agent        `json:"suppliedBy,omitempty"`
	SupportLevel   []SupportType `json:"supportLevel,omitempty"`
	ValidUntilTime *time.Time    `json:"validUntilTime,omitempty"`
}

func (a *Artifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

type Person struct {
	Agent
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
	Type                   string                 `json:"type"` //  "ExternalIdentifier",
	ExternalIdentifierType ExternalIdentifierType `json:"externalIdentifierType"`
	IssuingAuthority       string                 `json:"issuingAuthority,omitempty"`
	Identifier             string                 `json:"identifier"`
	IdentifierLocator      []string               `json:"identifierLocator,omitempty"`
}

// ExternalRef represents a reference to a resource outside SPDX-3.0 content
type ExternalRef struct {
	Comment         string          `json:"comment,omitempty"`
	ContentType     string          `json:"contentType,omitempty"`
	ExternalRefType ExternalRefType `json:"externalRefType,omitempty"`
	Locator         []string        `json:"locator,omitempty"`
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

// Hash represents a cryptographic hash value
type Hash struct {
	IntegrityMethod
	Algorithm HashAlgorithm `json:"algorithm"`
	HashValue string        `json:"hashValue"`
}

// PackageVerificationCode provides package integrity verification
type PackageVerificationCode struct {
	IntegrityMethod
	Algorithm                           HashAlgorithm `json:"algorithm"`
	HashValue                           string        `json:"hashValue"`
	PackageVerificationCodeExcludedFile []string      `json:"packageVerificationCodeExcludedFile,omitempty"`
}

// ExternalMap maps Element identifiers defined external to the SpdxDocument
type ExternalMap struct {
	DefiningArtifact string            `json:"definingArtifact,omitempty"`
	ExternalSpdxId   string            `json:"externalSpdxId"`
	LocationHint     string            `json:"locationHint,omitempty"`
	VerifiedUsing    []IntegrityMethod `json:"verifiedUsing,omitempty"`
}

// NamespaceMap allows shorter identifiers for namespace portions
type NamespaceMap struct {
	Namespace string `json:"namespace"`
	Prefix    string `json:"prefix"`
}

// Annotation provides additional information about elements
type Annotation struct {
	Node
	AnnotationType AnnotationType `json:"annotationType"`
	ContentType    string         `json:"contentType,omitempty"`
	Statement      string         `json:"statement,omitempty"`
	Subject        string         `json:"subject"`
}

func (a *Annotation) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

// SoftwareAgent represents a software program acting on a system
type SoftwareAgent struct {
	Agent
}

func (sa *SoftwareAgent) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sa, &sa.PreNode)
}

// Tool represents hardware/software used to carry out a function
type Tool struct {
	Node
}

func (t *Tool) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, t, &t.PreNode)
}

// LifecycleScopedRelationship parameterizes context for relationships
type LifecycleScopedRelationship struct {
	Relationship
	Scope LifecycleScopeType `json:"scope,omitempty"`
}

func (lsr *LifecycleScopedRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, lsr, &lsr.PreNode)
}

// ElementCollection is an abstract collection of Elements
type ElementCollection struct {
	Node
	Element            []ElementDescendant     `json:"element,omitempty"`
	RootElement        []ElementDescendant     `json:"rootElement,omitempty"`
	ProfileConformance []ProfileIdentifierType `json:"profileConformance,omitempty"`
}

// Bundle is a collection of Elements with shared context
type Bundle struct {
	ElementCollection
	Context string `json:"context,omitempty"`
}

func (b *Bundle) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}

type Bom struct {
	Bundle
	types.RootedNode
}

func (b *Bom) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}

type SpdxDocument struct {
	Bundle
	types.RootedNode
	DataLicense  Element        `json:"dataLicense,omitempty"` // This is simple licenseinfo but we can't loop
	ExternalMap  []ExternalMap  `json:"externalMap,omitempty"`
	NamespaceMap []NamespaceMap `json:"namespaceMap,omitempty"`
}

func (sd *SpdxDocument) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sd, &sd.PreNode)
}
