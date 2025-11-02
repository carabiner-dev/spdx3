package simplelicensing

import (
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const Prefix = "simplelicensing "

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"AnyLicenseInfo":                      &AnyLicenseInfo{},
		"LicenseExpression":                   &LicenseExpression{},
		"SimpleLicensingText":                 &SimpleLicensingText{},
		"simplelicensing_AnyLicenseInfo":      &AnyLicenseInfo{},
		"simplelicensing_LicenseExpression":   &LicenseExpression{},
		"simplelicensing_SimpleLicensingText": &SimpleLicensingText{},
	},
}

// AnyLicenseInfo is the abstract base class for license information
type AnyLicenseInfo struct {
	core.Element
}

func (ali *AnyLicenseInfo) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ali, &ali.PreNode)
}

// LicenseExpression represents licenses expressed as license expression strings
type LicenseExpression struct {
	AnyLicenseInfo
	CustomIdToUri      []core.DictionaryEntry `json:"customIdToUri,omitempty"`
	LicenseExpression  string                 `json:"simplelicensing_licenseExpression"`
	LicenseListVersion string                 `json:"simplelicensing_licenseListVersion,omitempty"`
}

func (le *LicenseExpression) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, le, &le.PreNode)
}

// SimpleLicensingText represents license text not on the SPDX License List
type SimpleLicensingText struct {
	core.Element
	LicenseText string `json:"licenseText"`
}

func (slt *SimpleLicensingText) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, slt, &slt.PreNode)
}
