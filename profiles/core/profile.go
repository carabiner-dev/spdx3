// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                           = "core"
	ElementClass                     = "Element"
	CreationInfoClass                = "CreationInfo"
	RelationshipClass                = "Relationship"
	AgentClass                       = "Agent"
	ArtifactClass                    = "Artifact"
	PersonClass                      = "Person"
	OrganizationClass                = "Organization"
	AnnotationClass                  = "Annotation"
	SoftwareAgentClass               = "SoftwareAgent"
	ToolClass                        = "Tool"
	LifecycleScopedRelationshipClass = "LifecycleScopedRelationship"
	ElementCollectionClass           = "ElementCollection"
	BundleClass                      = "Bundle"
	BomClass                         = "Bom"
	SpdxDocumentClass                = "SpdxDocument"
	HashClass                        = "Hash"
	PackageVerificationCodeClass     = "PackageVerificationCode"
	IndividualElementClass           = "IndividualElement"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		ElementClass:                     &Node{},
		CreationInfoClass:                &CreationInfo{},
		RelationshipClass:                &Relationship{},
		AgentClass:                       &Agent{},
		ArtifactClass:                    &Artifact{},
		PersonClass:                      &Person{},
		OrganizationClass:                &Organization{},
		AnnotationClass:                  &Annotation{},
		SoftwareAgentClass:               &SoftwareAgent{},
		ToolClass:                        &Tool{},
		LifecycleScopedRelationshipClass: &LifecycleScopedRelationship{},
		ElementCollectionClass:           &ElementCollection{},
		BundleClass:                      &Bundle{},
		BomClass:                         &Bom{},
		SpdxDocumentClass:                &SpdxDocument{},
		HashClass:                        &Hash{},
		PackageVerificationCodeClass:     &PackageVerificationCode{},
		IndividualElementClass:           &IndividualElement{},
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

func (r *Relationship) GetType() string {
	return RelationshipClass
}

func (r *Relationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, r, &r.PreNode)
}

func (r *Relationship) FromRelationship() {}
func (r *Relationship) FromElement()      {}

// Element is the base element, the node of the SPDX graph
type Element struct {
	Node
}

func (e *Element) GetType() string {
	return ElementClass
}

func (e *Element) FromElement() {}

// ElementDescendant is implemented by every class the model derives from
// Element. The marker sits on each direct subclass of Element rather than
// on Node, because not everything carrying Node's properties is an Element:
// Extension is a standalone class in the model.
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

func (a *Agent) GetType() string {
	return AgentClass
}

func (a *Agent) FromAgent()   {}
func (a *Agent) FromElement() {}

func (a *Agent) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

// Artifact represents a distinct article or unit within the digital domain
type Artifact struct {
	Node
	BuiltTime      *time.Time        `json:"builtTime,omitempty"`
	OriginatedBy   []AgentDescendant `json:"originatedBy,omitempty"`
	ReleaseTime    *time.Time        `json:"releaseTime,omitempty"`
	StandardName   []string          `json:"standardName,omitempty"`
	SuppliedBy     AgentDescendant   `json:"suppliedBy,omitempty"`
	SupportLevel   []SupportType     `json:"supportLevel,omitempty"`
	ValidUntilTime *time.Time        `json:"validUntilTime,omitempty"`
}

func (a *Artifact) GetType() string {
	return ArtifactClass
}

func (a *Artifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

func (a *Artifact) FromElement() {}

type Person struct {
	Agent
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier,omitempty"`
}

func (p *Person) GetType() string {
	return PersonClass
}

func (p *Person) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, p, &p.PreNode)
}

type Organization struct {
	Agent
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier,omitempty"`
}

func (o *Organization) GetType() string {
	return OrganizationClass
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
	Type            string          `json:"type,omitempty"`
	Comment         string          `json:"comment,omitempty"`
	ContentType     string          `json:"contentType,omitempty"`
	ExternalRefType ExternalRefType `json:"externalRefType,omitempty"`
	Locator         []string        `json:"locator,omitempty"`
}

// DictionaryEntry represents a key-value pair mapping
type DictionaryEntry struct {
	Type  string `json:"type,omitempty"`
	Key   string `json:"key"`
	Value string `json:"value,omitempty"`
}

// IntegrityMethod provides an independently reproducible mechanism for verification
type IntegrityMethod struct {
	base.PreNode
	Comment string `json:"comment,omitempty"`
}

func (im *IntegrityMethod) GetName() string             { return "" }
func (im *IntegrityMethod) GetCreationInfo() types.Node { return nil }
func (im *IntegrityMethod) FromIntegrityMethod()        {}

// IntegrityMethodDescendant is implemented by all integrity method types so
// that fields like verifiedUsing can dispatch them polymorphically.
type IntegrityMethodDescendant interface {
	types.Node
	FromIntegrityMethod()
}

// PositiveIntegerRange represents a tuple of two positive integers defining a range
type PositiveIntegerRange struct {
	Type              string `json:"type,omitempty"`
	BeginIntegerRange int    `json:"beginIntegerRange"`
	EndIntegerRange   int    `json:"endIntegerRange"`
}

// Hash represents a cryptographic hash value
type Hash struct {
	IntegrityMethod
	Algorithm HashAlgorithm `json:"algorithm"`
	HashValue string        `json:"hashValue"`
}

func (h *Hash) GetType() string {
	return HashClass
}

// PackageVerificationCode provides package integrity verification
type PackageVerificationCode struct {
	IntegrityMethod
	Algorithm                           HashAlgorithm `json:"algorithm"`
	HashValue                           string        `json:"hashValue"`
	PackageVerificationCodeExcludedFile []string      `json:"packageVerificationCodeExcludedFile,omitempty"`
}

func (pvc *PackageVerificationCode) GetType() string {
	return PackageVerificationCodeClass
}

// ExternalMap maps Element identifiers defined external to the SpdxDocument
type ExternalMap struct {
	base.PreNode
	DefiningArtifact string                      `json:"definingArtifact,omitempty"`
	ExternalSpdxId   string                      `json:"externalSpdxId"`
	LocationHint     string                      `json:"locationHint,omitempty"`
	VerifiedUsing    []IntegrityMethodDescendant `json:"verifiedUsing,omitempty"`
}

func (em *ExternalMap) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, em, &em.PreNode)
}

// NamespaceMap allows shorter identifiers for namespace portions
type NamespaceMap struct {
	Type      string `json:"type,omitempty"`
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

func (a *Annotation) GetType() string {
	return AnnotationClass
}

func (a *Annotation) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}

func (a *Annotation) FromElement() {}

// SoftwareAgent represents a software program acting on a system
type SoftwareAgent struct {
	Agent
}

func (sa *SoftwareAgent) GetType() string {
	return SoftwareAgentClass
}

func (sa *SoftwareAgent) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sa, &sa.PreNode)
}

// Tool represents hardware/software used to carry out a function
type Tool struct {
	Node
}

func (t *Tool) GetType() string {
	return ToolClass
}

func (t *Tool) FromElement() {}

func (t *Tool) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, t, &t.PreNode)
}

// LifecycleScopedRelationship parameterizes context for relationships
type LifecycleScopedRelationship struct {
	Relationship
	Scope LifecycleScopeType `json:"scope,omitempty"`
}

func (lsr *LifecycleScopedRelationship) GetType() string {
	return LifecycleScopedRelationshipClass
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

func (ec *ElementCollection) GetType() string {
	return ElementCollectionClass
}

func (ec *ElementCollection) FromElement() {}

// Bundle is a collection of Elements with shared context
type Bundle struct {
	ElementCollection
	Context string `json:"context,omitempty"`
}

func (b *Bundle) GetType() string {
	return BundleClass
}

func (b *Bundle) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}

type Bom struct {
	Bundle
}

func (b *Bom) GetType() string {
	return BomClass
}

func (b *Bom) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}

type SpdxDocument struct {
	Bundle
	DataLicense  types.Node     `json:"dataLicense,omitempty"` // This is simple licenseinfo but we can't loop
	Import       []ExternalMap  `json:"import,omitempty"`
	NamespaceMap []NamespaceMap `json:"namespaceMap,omitempty"`
}

func (sd *SpdxDocument) GetType() string {
	return SpdxDocumentClass
}

func (sd *SpdxDocument) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sd, &sd.PreNode)
}
