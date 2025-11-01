package simplelicensing

import (
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "simplelicensing "

// simplelicensing_LicenseExpression
type LicenseExpression struct {
	core.Node
	LicenseExpression  string `json:"simplelicensing_licenseExpression"`
	LicenseListVersion string `json:"simplelicensing_licenseListVersion"`
}

func (le *LicenseExpression) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, le, &le.PreNode)
}
