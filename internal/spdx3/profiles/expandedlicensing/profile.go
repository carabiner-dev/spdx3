package expandedlicensing

import (
	"reflect"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/simplelicensing"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "expandedlicensing"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]reflect.Type{
		"ListedLicense":             reflect.TypeOf(&ListedLicense{}),
		"CustomLicense":             reflect.TypeOf(&CustomLicense{}),
		"OrLaterOperator":           reflect.TypeOf(&OrLaterOperator{}),
		"WithAdditionOperator":      reflect.TypeOf(&WithAdditionOperator{}),
		"ConjunctiveLicenseSet":     reflect.TypeOf(&ConjunctiveLicenseSet{}),
		"DisjunctiveLicenseSet":     reflect.TypeOf(&DisjunctiveLicenseSet{}),
		"IndividualLicensingInfo":   reflect.TypeOf(&IndividualLicensingInfo{}),
		"ListedLicenseException":    reflect.TypeOf(&ListedLicenseException{}),
		"CustomLicenseAddition":     reflect.TypeOf(&CustomLicenseAddition{}),
	},
}

// ExtendableLicense is an abstract base for licenses that can be extended
type ExtendableLicense struct {
	simplelicensing.AnyLicenseInfo
}

// License is an abstract class representing a license text
type License struct {
	ExtendableLicense
	LicenseText               string   `json:"licenseText"`
	IsDeprecatedLicenseId     bool     `json:"isDeprecatedLicenseId,omitempty"`
	IsFsfLibre                bool     `json:"isFsfLibre,omitempty"`
	IsOsiApproved             bool     `json:"isOsiApproved,omitempty"`
	LicenseXml                string   `json:"licenseXml,omitempty"`
	ObsoletedBy               string   `json:"obsoletedBy,omitempty"`
	SeeAlso                   []string `json:"seeAlso,omitempty"`
	StandardLicenseHeader     string   `json:"standardLicenseHeader,omitempty"`
	StandardLicenseTemplate   string   `json:"standardLicenseTemplate,omitempty"`
}

// ListedLicense references a license from the official SPDX list
type ListedLicense struct {
	License
	DeprecatedVersion string `json:"deprecatedVersion,omitempty"`
	ListVersionAdded  string `json:"listVersionAdded,omitempty"`
}

func (ll *ListedLicense) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ll, &ll.PreNode)
}

// CustomLicense defines a custom license not on the standard list
type CustomLicense struct {
	License
}

func (cl *CustomLicense) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cl, &cl.PreNode)
}

// OrLaterOperator represents the "or later version" operator
type OrLaterOperator struct {
	simplelicensing.AnyLicenseInfo
	SubjectLicense string `json:"subjectLicense"`
}

func (olo *OrLaterOperator) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, olo, &olo.PreNode)
}

// WithAdditionOperator represents the "with addition" operator
type WithAdditionOperator struct {
	simplelicensing.AnyLicenseInfo
	SubjectAddition         string `json:"subjectAddition"`
	SubjectExtendableLicense string `json:"subjectExtendableLicense"`
}

func (wao *WithAdditionOperator) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, wao, &wao.PreNode)
}

// ConjunctiveLicenseSet represents a conjunction of multiple licenses (AND)
type ConjunctiveLicenseSet struct {
	simplelicensing.AnyLicenseInfo
	Member []string `json:"member"`
}

func (cls *ConjunctiveLicenseSet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cls, &cls.PreNode)
}

// DisjunctiveLicenseSet represents a disjunction of multiple licenses (OR)
type DisjunctiveLicenseSet struct {
	simplelicensing.AnyLicenseInfo
	Member []string `json:"member"`
}

func (dls *DisjunctiveLicenseSet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, dls, &dls.PreNode)
}

// IndividualLicensingInfo captures specific licensing information
type IndividualLicensingInfo struct {
	simplelicensing.AnyLicenseInfo
}

func (ili *IndividualLicensingInfo) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ili, &ili.PreNode)
}

// LicenseAddition is an abstract class for additions/exceptions to licenses
type LicenseAddition struct {
	core.Node
	AdditionText              string   `json:"additionText"`
	IsDeprecatedAdditionId    bool     `json:"isDeprecatedAdditionId,omitempty"`
	LicenseXml                string   `json:"licenseXml,omitempty"`
	ObsoletedBy               string   `json:"obsoletedBy,omitempty"`
	SeeAlso                   []string `json:"seeAlso,omitempty"`
	StandardAdditionTemplate  string   `json:"standardAdditionTemplate,omitempty"`
}

// ListedLicenseException references an exception from the SPDX list
type ListedLicenseException struct {
	LicenseAddition
	DeprecatedVersion string `json:"deprecatedVersion,omitempty"`
	ListVersionAdded  string `json:"listVersionAdded,omitempty"`
}

func (lle *ListedLicenseException) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, lle, &lle.PreNode)
}

// CustomLicenseAddition represents a custom addition to a license
type CustomLicenseAddition struct {
	LicenseAddition
}

func (cla *CustomLicenseAddition) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, cla, &cla.PreNode)
}
