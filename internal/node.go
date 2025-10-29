package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

var ErrUnsupportedNodeType = errors.New("unsupported node type")

type Renderer struct{}

func (r *Renderer) Render(env *Envelope, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(*env)
}

type Parser struct {
}

func (p *Parser) Parse(r io.Reader) (*Envelope, error) {
	dec := json.NewDecoder(r)
	env := &Envelope{}
	if err := dec.Decode(env); err != nil {
		return nil, err
	}
	return env, nil
}

// Envelope
type Envelope struct {
	Context string `json:"@context"`
	Graph   Graph  `json:"@graph"`
}

type Graph []Node

func (g *Graph) UnmarshalJSON(data []byte) error {
	list := []json.RawMessage{}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshaling graph: %w", err)
	}

	for i, prenodeData := range list {
		// Parse the entry to a prenode
		var prenode = &PreNode{}
		if err := json.Unmarshal(prenodeData, prenode); err != nil {
			return fmt.Errorf("parsing node #%d: %w", i, err)
		}

		var n Node
		switch prenode.Type {
		case "CreationInfo":
			n = &CreationInfo{}
		case "Person":
			n = &Person{}
		case "Organization":
			n = &Organization{}
		case "SpdxDocument":
			n = &SpdxDocument{}
		case "Bom":
			n = &Bom{}
		case "dataset_DatasetPackage":
			n = &DatasetPackage{}
		case "software_File":
			n = &File{}
		case "Relationship":
			n = &Relationship{}
		case "simplelicensing_LicenseExpression":
			n = &LicenseExpression{}
		default:
			return fmt.Errorf("parsing type %q: %w", prenode.Type, ErrUnsupportedNodeType)
		}
		if err := json.Unmarshal(prenodeData, n); err != nil {
			return err
		}
		// Dedupe IDs
		*g = append(*g, n)
	}

	return nil
}

type Node interface {
	GetSPDXID() string
	GetID() string
	GetType() string
	GetName() string
	GetCreationInfo() Marshable
}

type PreNode struct {
	SPDXID string `json:"spdxID"`
	ID     string `json:"@id"`
	Type   string `json:"type"`
}

func (pn *PreNode) GetSPDXID() string {
	return pn.SPDXID
}

func (pn *PreNode) GetID() string {
	return pn.ID
}

func (pn *PreNode) GetType() string {
	return pn.Type
}

type BaseNode struct {
	PreNode
	Name         string        `json:"name"`
	CreationInfo *CreationInfo `json:"creationInfo"`
	Comment      string        `json:"comment"`
	Description  string        `json:"description"`
}

func (bn *BaseNode) GetCreationInfo() Marshable {
	return bn.CreationInfo
}

func (bn *BaseNode) GetName() string {
	return bn.Name
}

type Marshable interface {
	MarshalJSON() ([]byte, error)
}

type SoftwareNode struct {
	DownloadLocation string `json:"software_downloadLocation"`
	PrimaryPurpose   string `json:"software_primaryPurpose"`
}

type DataSetNode struct {
	BuiltTime                       *time.Time `json:"builtTime"`
	ReleaseTime                     *time.Time `json:"releaseTime"`
	ConfidentialityLevel            string     `json:"confidentialityLevel"`
	DataPreprocessing               []string   `json:"dataset_dataPreprocessing"`
	DatasetAvailability             string     `json:"dataset_datasetAvailability"`
	DataCollectionProcess           string     `json:"dataset_dataCollectionProcess"`
	DatasetSize                     int64      `json:"dataset_datasetSize"`
	DatasetType                     []string   `json:"dataset_datasetType"`
	DatasetUpdateMechanism          string     `json:"dataset_datasetUpdateMechanism"`
	HasSensitivePersonalInformation string     `json:"dataset_hasSensitivePersonalInformation"`
	IntendedUse                     string     `json:"dataset_intendedUse"`
	KnownBias                       []string   `json:"dataset_knownBias"`
}

type DatasetPackage struct {
	BaseNode
	DataSetNode
}

type File struct {
	BaseNode
	SoftwareNode
	BuiltTime    *time.Time `json:"builtTime"`
	ReleaseTime  *time.Time `json:"releaseTime"`
	OriginatedBy []string   `json:"originatedBy"`
}

type Relationship struct {
	BaseNode
	From              string
	To                []string
	RelationshipTypes string `json:"relationshipType"`
}

// simplelicensing_LicenseExpression
type LicenseExpression struct {
	BaseNode
	LicenseExpression  string `json:"simplelicensing_licenseExpression"`
	LicenseListVersion string `json:"simplelicensing_licenseListVersion"`
}

type Bom struct {
	BaseNode
	RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

type SpdxDocument struct {
	BaseNode
	RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

type Organization struct {
	BaseNode
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

type Person struct {
	BaseNode
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

type ExternalIdentifier struct {
	Type                   string   `json:"type"` //  "ExternalIdentifier",
	ExternalIdentifierType string   `json:"externalIdentifierType"`
	IssuingAuthority       string   `json:"issuingAuthority"`
	Identifier             string   `json:"identifier"`
	IdentifierLocator      []string `json:"identifierLocator"`
}

type RootedNode struct {
	RootElement []string `json:"rotElement"`
}

type CreationInfo struct {
	PreNode
	RootedNode
	root        bool       `json:""`
	Name        string     `json:"name"`
	SpecVersion string     `json:"specVersion"`
	CreatedBy   []string   `json:"createdBy"`
	Created     *time.Time `json:"created"`
}

func (ci *CreationInfo) GetName() string {
	return ci.Name
}

func (ci *CreationInfo) GetCreationInfo() Marshable {
	return ci
}

func (ci *CreationInfo) MarshalJSON() ([]byte, error) {
	if ci.root {
		return json.Marshal(ci)
	}

	return fmt.Appendf(nil, "\"_:%s\"", ci.ID), nil
}

func (ci *CreationInfo) UnmarshalJSON(data []byte) error {
	type Alias CreationInfo
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ci),
	}
	if err := json.Unmarshal(data, aux); err == nil {
		return nil
	}

	var s string
	if errs := json.Unmarshal(data, &s); errs == nil {
		ci.ID = strings.TrimPrefix(s, "_:")
		return nil
	}

	return nil
}
