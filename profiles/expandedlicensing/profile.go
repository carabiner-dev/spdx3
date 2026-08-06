// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package expandedlicensing

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/simplelicensing"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                      = "expandedlicensing"
	ListedLicenseType           = "ListedLicense"
	CustomLicenseType           = "CustomLicense"
	OrLaterOperatorType         = "OrLaterOperator"
	WithAdditionOperatorType    = "WithAdditionOperator"
	ConjunctiveLicenseSetType   = "ConjunctiveLicenseSet"
	DisjunctiveLicenseSetType   = "DisjunctiveLicenseSet"
	IndividualLicensingInfoType = "IndividualLicensingInfo"
	ListedLicenseExceptionType  = "ListedLicenseException"
	CustomLicenseAdditionType   = "CustomLicenseAddition"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		ListedLicenseType:                                         &ListedLicense{},
		CustomLicenseType:                                         &CustomLicense{},
		OrLaterOperatorType:                                       &OrLaterOperator{},
		WithAdditionOperatorType:                                  &WithAdditionOperator{},
		ConjunctiveLicenseSetType:                                 &ConjunctiveLicenseSet{},
		DisjunctiveLicenseSetType:                                 &DisjunctiveLicenseSet{},
		IndividualLicensingInfoType:                               &IndividualLicensingInfo{},
		ListedLicenseExceptionType:                                &ListedLicenseException{},
		CustomLicenseAdditionType:                                 &CustomLicenseAddition{},
		fmt.Sprintf("%s_%s", Prefix, ListedLicenseType):           &ListedLicense{},
		fmt.Sprintf("%s_%s", Prefix, CustomLicenseType):           &CustomLicense{},
		fmt.Sprintf("%s_%s", Prefix, OrLaterOperatorType):         &OrLaterOperator{},
		fmt.Sprintf("%s_%s", Prefix, WithAdditionOperatorType):    &WithAdditionOperator{},
		fmt.Sprintf("%s_%s", Prefix, ConjunctiveLicenseSetType):   &ConjunctiveLicenseSet{},
		fmt.Sprintf("%s_%s", Prefix, DisjunctiveLicenseSetType):   &DisjunctiveLicenseSet{},
		fmt.Sprintf("%s_%s", Prefix, IndividualLicensingInfoType): &IndividualLicensingInfo{},
		fmt.Sprintf("%s_%s", Prefix, ListedLicenseExceptionType):  &ListedLicenseException{},
		fmt.Sprintf("%s_%s", Prefix, CustomLicenseAdditionType):   &CustomLicenseAddition{},
	},
}

// ExtendableLicense is an abstract base for licenses that can be extended
type ExtendableLicense struct {
	simplelicensing.AnyLicenseInfo
}

// License is an abstract class representing a license text
type License struct {
	ExtendableLicense
	LicenseText             string   `json:"simplelicensing_licenseText"`
	IsDeprecatedLicenseId   bool     `json:"expandedlicensing_isDeprecatedLicenseId,omitempty"`
	IsFsfLibre              bool     `json:"expandedlicensing_isFsfLibre,omitempty"`
	IsOsiApproved           bool     `json:"expandedlicensing_isOsiApproved,omitempty"`
	LicenseXml              string   `json:"expandedlicensing_licenseXml,omitempty"`
	ObsoletedBy             string   `json:"expandedlicensing_obsoletedBy,omitempty"`
	SeeAlso                 []string `json:"expandedlicensing_seeAlso,omitempty"`
	StandardLicenseHeader   string   `json:"expandedlicensing_standardLicenseHeader,omitempty"`
	StandardLicenseTemplate string   `json:"expandedlicensing_standardLicenseTemplate,omitempty"`
}

// ListedLicense references a license from the official SPDX list
type ListedLicense struct {
	License
	DeprecatedVersion string `json:"expandedlicensing_deprecatedVersion,omitempty"`
	ListVersionAdded  string `json:"expandedlicensing_listVersionAdded,omitempty"`
}

func (ll *ListedLicense) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, ListedLicenseType)
}

func (ll *ListedLicense) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ll, &ll.PreNode)
}

// CustomLicense defines a custom license not on the standard list
type CustomLicense struct {
	License
}

func (cl *CustomLicense) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CustomLicenseType)
}

func (cl *CustomLicense) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cl, &cl.PreNode)
}

// OrLaterOperator represents the "or later version" operator
type OrLaterOperator struct {
	simplelicensing.AnyLicenseInfo
	SubjectLicense string `json:"expandedlicensing_subjectLicense"`
}

func (olo *OrLaterOperator) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, OrLaterOperatorType)
}

func (olo *OrLaterOperator) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, olo, &olo.PreNode)
}

// WithAdditionOperator represents the "with addition" operator
type WithAdditionOperator struct {
	simplelicensing.AnyLicenseInfo
	SubjectAddition          string `json:"expandedlicensing_subjectAddition"`
	SubjectExtendableLicense string `json:"expandedlicensing_subjectExtendableLicense"`
}

func (wao *WithAdditionOperator) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, WithAdditionOperatorType)
}

func (wao *WithAdditionOperator) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, wao, &wao.PreNode)
}

// ConjunctiveLicenseSet represents a conjunction of multiple licenses (AND)
type ConjunctiveLicenseSet struct {
	simplelicensing.AnyLicenseInfo
	Member []string `json:"expandedlicensing_member"`
}

func (cls *ConjunctiveLicenseSet) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, ConjunctiveLicenseSetType)
}

func (cls *ConjunctiveLicenseSet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cls, &cls.PreNode)
}

// DisjunctiveLicenseSet represents a disjunction of multiple licenses (OR)
type DisjunctiveLicenseSet struct {
	simplelicensing.AnyLicenseInfo
	Member []string `json:"expandedlicensing_member"`
}

func (dls *DisjunctiveLicenseSet) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, DisjunctiveLicenseSetType)
}

func (dls *DisjunctiveLicenseSet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, dls, &dls.PreNode)
}

// IndividualLicensingInfo captures specific licensing information
type IndividualLicensingInfo struct {
	simplelicensing.AnyLicenseInfo
}

func (ili *IndividualLicensingInfo) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, IndividualLicensingInfoType)
}

func (ili *IndividualLicensingInfo) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ili, &ili.PreNode)
}

// LicenseAddition is an abstract class for additions/exceptions to licenses
type LicenseAddition struct {
	core.Node
	AdditionText             string   `json:"expandedlicensing_additionText"`
	IsDeprecatedAdditionId   bool     `json:"expandedlicensing_isDeprecatedAdditionId,omitempty"`
	LicenseXml               string   `json:"expandedlicensing_licenseXml,omitempty"`
	ObsoletedBy              string   `json:"expandedlicensing_obsoletedBy,omitempty"`
	SeeAlso                  []string `json:"expandedlicensing_seeAlso,omitempty"`
	StandardAdditionTemplate string   `json:"expandedlicensing_standardAdditionTemplate,omitempty"`
}

// ListedLicenseException references an exception from the SPDX list
type ListedLicenseException struct {
	LicenseAddition
	DeprecatedVersion string `json:"expandedlicensing_deprecatedVersion,omitempty"`
	ListVersionAdded  string `json:"expandedlicensing_listVersionAdded,omitempty"`
}

func (lle *ListedLicenseException) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, ListedLicenseExceptionType)
}

func (lle *ListedLicenseException) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, lle, &lle.PreNode)
}

// CustomLicenseAddition represents a custom addition to a license
type CustomLicenseAddition struct {
	LicenseAddition
}

func (cla *CustomLicenseAddition) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CustomLicenseAdditionType)
}

func (cla *CustomLicenseAddition) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cla, &cla.PreNode)
}
