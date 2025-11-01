package dispatch

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

func New() types.Dispatcher {
	return &Default{}
}

// Default is the default type dispatcher. Its job is to
// detect its type and pass the correct node type to the
// json unmarshaller.
type Default struct{}

// UnmarshalNode takes raw JSON data and unmarshals it into the appropriate
// concrete type based on the "type" field. This is used to handle polymorphic
// node types in SPDX3 JSON-LD documents.
func (d *Default) UnmarshalNode(prenodeData []byte) (types.Node, error) {
	var prenode = &base.PreNode{}
	// Parse the entry to a prenode to determine its type
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
		return nil, fmt.Errorf("parsing type %q: %w", prenode.Type, types.ErrUnsupportedNodeType)
	}
	if err := json.Unmarshal(prenodeData, n); err != nil {
		return nil, err
	}
	return n, nil
}
