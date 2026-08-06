// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package expandedlicensing

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/simplelicensing"
)

// IRIs of the licensing individuals the spec predefines. Note that they
// live in the Licensing namespace, not ExpandedLicensing, and use the short
// names None and NoAssertion. Documents reference the individuals through
// these IRIs rather than serializing them, so license references should be
// matched against these constants. The IRIs are pinned to the 3.0.1
// namespace and re-version with the spec.
const (
	NoneLicenseIRI        = "https://spdx.org/rdf/3.0.1/terms/Licensing/None"
	NoAssertionLicenseIRI = "https://spdx.org/rdf/3.0.1/terms/Licensing/NoAssertion"
)

// IndividualIRIs returns the IRIs of the licensing individuals the spec
// predefines. Documents reference them without ever serializing them, so
// they are legitimate reference targets that no graph contains.
func IndividualIRIs() []string {
	return []string{
		NoneLicenseIRI,
		NoAssertionLicenseIRI,
	}
}

// NoneLicense is the individual asserting that the SPDX data creator
// determined no license is present. Shared read-only instance; compare
// references by IRI.
var NoneLicense = newLicensingIndividual(NoneLicenseIRI, "NONE")

// NoAssertionLicense is the individual asserting that the SPDX data creator
// makes no claim about license information. Shared read-only instance;
// compare references by IRI.
var NoAssertionLicense = newLicensingIndividual(NoAssertionLicenseIRI, "NOASSERTION")

func newLicensingIndividual(iri, name string) *IndividualLicensingInfo {
	return &IndividualLicensingInfo{
		AnyLicenseInfo: simplelicensing.AnyLicenseInfo{
			Element: core.Element{
				Node: core.Node{
					PreNode: base.PreNode{
						SPDXID: iri,
						Type:   fmt.Sprintf("%s_%s", Prefix, IndividualLicensingInfoType),
					},
					Name: name,
				},
			},
		},
	}
}
