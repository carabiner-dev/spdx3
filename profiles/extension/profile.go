package extension

import (
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "extension"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"CdxPropertiesExtension": &CdxPropertiesExtension{},
	},
}

// CdxPropertyEntry represents a name-value pair for CycloneDX compatible properties
type CdxPropertyEntry struct {
	CdxPropName  string `json:"cdxPropName"`
	CdxPropValue string `json:"cdxPropValue,omitempty"`
}

// CdxPropertiesExtension provides CycloneDX-compatible property extensions
type CdxPropertiesExtension struct {
	core.Extension
	CdxProperty []CdxPropertyEntry `json:"cdxProperty"`
}

func (cpe *CdxPropertiesExtension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cpe, &cpe.PreNode)
}
