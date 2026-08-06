// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

// IndividualElement is the concrete Element subclass that types the named
// individuals the Core profile defines (NoneElement, NoAssertionElement).
type IndividualElement struct {
	Element
}

func (ie *IndividualElement) GetType() string {
	return IndividualElementClass
}

func (ie *IndividualElement) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ie, &ie.PreNode)
}

// IRIs of the individuals the Core profile predefines. Documents reference
// the individuals through these IRIs rather than serializing them, so
// element references should be matched against these constants. The IRIs
// are pinned to the 3.0.1 namespace and re-version with the spec.
const (
	NoneElementIRI        = "https://spdx.org/rdf/3.0.1/terms/Core/NoneElement"
	NoAssertionElementIRI = "https://spdx.org/rdf/3.0.1/terms/Core/NoAssertionElement"
	SpdxOrganizationIRI   = "https://spdx.org/rdf/3.0.1/terms/Core/SpdxOrganization"
)

// IndividualIRIs returns the IRIs of the individuals the Core profile
// predefines. Documents reference them without ever serializing them, so
// they are legitimate reference targets that no graph contains.
func IndividualIRIs() []string {
	return []string{
		NoneElementIRI,
		NoAssertionElementIRI,
		SpdxOrganizationIRI,
	}
}

// NoneElement is the individual asserting that no element exists: it
// represents a set of Elements with cardinality zero. Shared read-only
// instance; compare references by IRI.
var NoneElement = &IndividualElement{
	Element: Element{
		Node: Node{
			PreNode: base.PreNode{
				SPDXID: NoneElementIRI,
				Type:   IndividualElementClass,
			},
			Name: "NONE",
		},
	},
}

// NoAssertionElement is the individual asserting that the document creator
// makes no claim about the element's value. Shared read-only instance;
// compare references by IRI.
var NoAssertionElement = &IndividualElement{
	Element: Element{
		Node: Node{
			PreNode: base.PreNode{
				SPDXID: NoAssertionElementIRI,
				Type:   IndividualElementClass,
			},
			Name: "NOASSERTION",
		},
	},
}

// SpdxOrganization is the individual representing the SPDX Project itself,
// used as the creator of spec-defined content. Shared read-only instance;
// compare references by IRI.
var SpdxOrganization = &Organization{
	Agent: Agent{
		Node: Node{
			PreNode: base.PreNode{
				SPDXID: SpdxOrganizationIRI,
				Type:   OrganizationClass,
			},
			Name: "SPDX Project",
		},
	},
}
