// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package extension

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                     = "extension"
	CdxPropertiesExtensionType = "CdxPropertiesExtension"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		CdxPropertiesExtensionType:                               &CdxPropertiesExtension{},
		fmt.Sprintf("%s_%s", Prefix, CdxPropertiesExtensionType): &CdxPropertiesExtension{},
	},
}

// CdxPropertyEntry represents a name-value pair for CycloneDX compatible properties
type CdxPropertyEntry struct {
	CdxPropName  string `json:"extension_cdxPropName"`
	CdxPropValue string `json:"extension_cdxPropValue,omitempty"`
}

// CdxPropertiesExtension provides CycloneDX-compatible property extensions
type CdxPropertiesExtension struct {
	core.Extension
	CdxProperty []CdxPropertyEntry `json:"extension_cdxProperty"`
}

func (cpe *CdxPropertiesExtension) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CdxPropertiesExtensionType)
}

func (cpe *CdxPropertiesExtension) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cpe, &cpe.PreNode)
}
