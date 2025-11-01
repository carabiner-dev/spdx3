package spdx3

import (
	"encoding/json"
	"fmt"

	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/dataset"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/simplelicensing"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/software"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

func NewNodeUnarshaler() *unmarshal.Node {
	return &unmarshal.Node{
		DispatchFn: UnmarshalNode,
	}
}

// Envelope
type Envelope struct {
	Context string `json:"@context"`
	Graph   Graph  `json:"@graph"`
}

type Graph []types.Node

// UnmarshalNode takes raw JSON data and unmarshals it into the appropriate
// concrete type based on the "type" field. This is used to handle polymorphic
// node types in SPDX3 JSON-LD documents.
func UnmarshalNode(prenodeData []byte) (types.Node, error) {
	// Parse the entry to a prenode to determine its type
	var prenode = &base.PreNode{}
	if err := json.Unmarshal(prenodeData, prenode); err != nil {
		return nil, fmt.Errorf("parsing node: %w", err)
	}

	var n types.Node
	switch prenode.Type {
	case "CreationInfo":
		n = &core.CreationInfo{}
	case "Person":
		n = &core.Person{}
	case "Organization":
		n = &core.Organization{}
	case "SpdxDocument":
		n = &core.SpdxDocument{}
	case "Bom":
		n = &core.Bom{}
	case "dataset_DatasetPackage":
		n = &dataset.Package{}
	case "software_File":
		n = &software.File{}
	case "Relationship":
		n = &core.Relationship{}
	case "simplelicensing_LicenseExpression":
		n = &simplelicensing.LicenseExpression{}
	default:
		return nil, fmt.Errorf("parsing type %q: %w", prenode.Type, ErrUnsupportedNodeType)
	}
	if err := json.Unmarshal(prenodeData, n); err != nil {
		return nil, err
	}
	return n, nil
}

// UnmarshalJSON unmarshalls the JSONLD graph into nodes typed to their kinds.
func (g *Graph) UnmarshalJSON(data []byte) error {
	list := []json.RawMessage{}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshaling graph: %w", err)
	}

	for i, prenodeData := range list {
		n, err := UnmarshalNode(prenodeData)
		if err != nil {
			return fmt.Errorf("unmarshaling node #%d: %w", i, err)
		}
		// TODO(puerco): Dedupe IDs of the resulting node
		//  TODO(puerco): Dedupe IDs of any fields that are nodes by looking up
		// exising nodes with the same ID
		*g = append(*g, n)
	}

	return nil
}
