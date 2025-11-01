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
)

// Envelope
type Envelope struct {
	Context string `json:"@context"`
	Graph   Graph  `json:"@graph"`
}

type Graph []types.Node

// UnmarshalJSON unmarshalls the JSONLD graph into nodes typed to their kinds.
func (g *Graph) UnmarshalJSON(data []byte) error {
	list := []json.RawMessage{}
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("unmarshaling graph: %w", err)
	}

	for i, prenodeData := range list {
		// Parse the entry to a prenode to determine its type
		var prenode = &base.PreNode{}
		if err := json.Unmarshal(prenodeData, prenode); err != nil {
			return fmt.Errorf("parsing node #%d: %w", i, err)
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
			return fmt.Errorf("parsing type %q: %w", prenode.Type, ErrUnsupportedNodeType)
		}
		if err := json.Unmarshal(prenodeData, n); err != nil {
			return err
		}
		// TODO(puerco): Dedupe IDs
		*g = append(*g, n)
	}

	return nil
}
