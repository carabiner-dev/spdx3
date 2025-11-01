package extension

import (
	"reflect"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "extension"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]reflect.Type{
		"Extension":              reflect.TypeOf(&Extension{}),
		"CdxPropertiesExtension": reflect.TypeOf(&CdxPropertiesExtension{}),
	},
}

// Extension is the abstract base class for all SPDX extensions
type Extension struct {
	core.Node
}

func (e *Extension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}

// CdxPropertyEntry represents a name-value pair for CycloneDX compatible properties
type CdxPropertyEntry struct {
	CdxPropName  string `json:"cdxPropName"`
	CdxPropValue string `json:"cdxPropValue,omitempty"`
}

// CdxPropertiesExtension provides CycloneDX-compatible property extensions
type CdxPropertiesExtension struct {
	Extension
	CdxProperty []CdxPropertyEntry `json:"cdxProperty"`
}

func (cpe *CdxPropertiesExtension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cpe, &cpe.PreNode)
}
