package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
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

// UnmarshalJSON unmarshalls the JSONLD graph into nodes typed to their kinds.
func (g *Graph) UnmarshalJSON(data []byte) error {
	list := []json.RawMessage{}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshaling graph: %w", err)
	}

	for i, prenodeData := range list {
		// Parse the entry to a prenode to determine its type
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

func (dp *DatasetPackage) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, dp, &dp.PreNode)
}

// File
type File struct {
	BaseNode
	SoftwareNode
	BuiltTime    *time.Time `json:"builtTime"`
	ReleaseTime  *time.Time `json:"releaseTime"`
	OriginatedBy []string   `json:"originatedBy"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, f, &f.PreNode)
}

type Relationship struct {
	BaseNode
	From              string
	To                []string
	RelationshipTypes string `json:"relationshipType"`
}

func (r *Relationship) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, r, &r.PreNode)
}

// simplelicensing_LicenseExpression
type LicenseExpression struct {
	BaseNode
	LicenseExpression  string `json:"simplelicensing_licenseExpression"`
	LicenseListVersion string `json:"simplelicensing_licenseListVersion"`
}

func (le *LicenseExpression) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, le, &le.PreNode)
}

type Bom struct {
	BaseNode
	RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

func (b *Bom) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, b, &b.PreNode)
}

type SpdxDocument struct {
	BaseNode
	RootedNode
	ProfileConformance []string `json:"profileConformance"`
}

func (sd *SpdxDocument) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, sd, &sd.PreNode)
}

type Organization struct {
	BaseNode
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

func (o *Organization) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, o, &o.PreNode)
}

type Person struct {
	BaseNode
	ExternalIdentifier []ExternalIdentifier `json:"externalIdentifier"`
}

func (p *Person) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, p, &p.PreNode)
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

// unmarshalNode is a universal unmarshaling helper for any type that embeds PreNode.
// It handles both full object serialization and string reference serialization (e.g., "_:id").
func unmarshalNode(data []byte, target interface{}, preNodePtr *PreNode) error {
	// First, check if it's a string reference
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		preNodePtr.ID = strings.TrimPrefix(s, "_:")
		return nil
	}

	// Otherwise, it's an object. Unmarshal into a map first to get all fields
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Use reflection to set fields directly
	return unmarshalFields(raw, target)
}

// unmarshalFields recursively unmarshals fields including embedded structs
func unmarshalFields(raw map[string]json.RawMessage, target interface{}) error {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if field.Anonymous {
			// Handle embedded structs by recursively unmarshaling
			if fieldValue.CanAddr() {
				if err := unmarshalFields(raw, fieldValue.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		// Get the JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// Handle json tag options like "name,omitempty"
		tagName := strings.Split(jsonTag, ",")[0]
		if tagName == "" {
			continue
		}

		// Check if we have data for this field
		if rawData, ok := raw[tagName]; ok {
			// Create a new value of the field's type
			newVal := reflect.New(fieldValue.Type())
			if err := json.Unmarshal(rawData, newVal.Interface()); err != nil {
				return err
			}
			fieldValue.Set(newVal.Elem())
		}
	}

	return nil
}

func (ci *CreationInfo) UnmarshalJSON(data []byte) error {
	return unmarshalNode(data, ci, &ci.PreNode)
}
