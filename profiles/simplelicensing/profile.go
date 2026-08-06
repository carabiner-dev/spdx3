// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package simplelicensing

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                   = "simplelicensing"
	AnyLicenseInfoType       = "AnyLicenseInfo"
	LicenseExpressionType    = "LicenseExpression"
	SimpleLicensingTextType  = "SimpleLicensingText"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		AnyLicenseInfoType:                                 &AnyLicenseInfo{},
		LicenseExpressionType:                              &LicenseExpression{},
		SimpleLicensingTextType:                            &SimpleLicensingText{},
		fmt.Sprintf("%s_%s", Prefix, AnyLicenseInfoType):  &AnyLicenseInfo{},
		fmt.Sprintf("%s_%s", Prefix, LicenseExpressionType): &LicenseExpression{},
		fmt.Sprintf("%s_%s", Prefix, SimpleLicensingTextType): &SimpleLicensingText{},
	},
}

// AnyLicenseInfo is the abstract base class for license information
type AnyLicenseInfo struct {
	core.Element
}

func (ali *AnyLicenseInfo) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, AnyLicenseInfoType)
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

func (le *LicenseExpression) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, LicenseExpressionType)
}

func (le *LicenseExpression) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, le, &le.PreNode)
}

// SimpleLicensingText represents license text not on the SPDX License List
type SimpleLicensingText struct {
	core.Element
	LicenseText string `json:"licenseText"`
}

func (slt *SimpleLicensingText) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, SimpleLicensingTextType)
}

func (slt *SimpleLicensingText) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, slt, &slt.PreNode)
}
