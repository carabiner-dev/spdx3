package simplelicensing

import (
	"reflect"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "simplelicensing "

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]reflect.Type{
		"AnyLicenseInfo":      reflect.TypeOf(&AnyLicenseInfo{}),
		"LicenseExpression":   reflect.TypeOf(&LicenseExpression{}),
		"SimpleLicensingText": reflect.TypeOf(&SimpleLicensingText{}),
	},
}

// AnyLicenseInfo is the abstract base class for license information
type AnyLicenseInfo struct {
	core.Node
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
	core.Node
	LicenseText string `json:"licenseText"`
}

func (slt *SimpleLicensingText) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, slt, &slt.PreNode)
}
