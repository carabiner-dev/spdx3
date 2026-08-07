// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/profiles/simplelicensing"
)

// "MIT AND (Apache-2.0 OR BSD-3-Clause)": a license set holding another
// license set. The members used to be []string, so any member written as an
// object aborted the parse and an expression like this was unrepresentable.
const nestedLicenseDoc = `{
	"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
	"@graph": [
		{
			"spdxId": "https://example.com/and",
			"type": "expandedlicensing_ConjunctiveLicenseSet",
			"expandedlicensing_member": [
				"https://example.com/mit",
				{
					"spdxId": "https://example.com/or",
					"type": "expandedlicensing_DisjunctiveLicenseSet",
					"expandedlicensing_member": [
						"https://example.com/apache",
						{
							"spdxId": "https://example.com/bsd",
							"type": "expandedlicensing_CustomLicense",
							"simplelicensing_licenseText": "the BSD text"
						}
					]
				}
			]
		}
	]
}`

func TestNestedLicenseSets(t *testing.T) {
	env, err := NewParser().Parse(strings.NewReader(nestedLicenseDoc))
	require.NoError(t, err)

	and, ok := env.Graph[0].(*expandedlicensing.ConjunctiveLicenseSet)
	require.True(t, ok)
	require.Len(t, and.Member, 2)

	// A member written as a reference stays a reference...
	require.Equal(t, "https://example.com/mit", and.Member[0].GetID())

	// ...and one written inline keeps its concrete type, recursively.
	or, ok := and.Member[1].(*expandedlicensing.DisjunctiveLicenseSet)
	require.True(t, ok)
	require.Len(t, or.Member, 2)
	require.Equal(t, "https://example.com/apache", or.Member[0].GetID())

	bsd, ok := or.Member[1].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.Equal(t, "the BSD text", bsd.LicenseText)

	// The whole expression survives a round trip.
	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))
	require.Contains(t, buf.String(), "the BSD text")

	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	reAnd, ok := reparsed.Graph[0].(*expandedlicensing.ConjunctiveLicenseSet)
	require.True(t, ok)
	reOr, ok := reAnd.Member[1].(*expandedlicensing.DisjunctiveLicenseSet)
	require.True(t, ok)
	reBsd, ok := reOr.Member[1].(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.Equal(t, "the BSD text", reBsd.LicenseText)
}

// The operators take the classes the model gives as their range, so an
// inline value in any of them dispatches to its concrete type too.
func TestLicenseOperators(t *testing.T) {
	doc := `{
		"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
		"@graph": [
			{
				"spdxId": "https://example.com/orlater",
				"type": "expandedlicensing_OrLaterOperator",
				"expandedlicensing_subjectLicense": {
					"spdxId": "https://example.com/gpl",
					"type": "expandedlicensing_CustomLicense",
					"simplelicensing_licenseText": "the GPL text"
				}
			},
			{
				"spdxId": "https://example.com/with",
				"type": "expandedlicensing_WithAdditionOperator",
				"expandedlicensing_subjectExtendableLicense": "https://example.com/gpl",
				"expandedlicensing_subjectAddition": {
					"spdxId": "https://example.com/exception",
					"type": "expandedlicensing_CustomLicenseAddition",
					"expandedlicensing_additionText": "the exception text"
				}
			}
		]
	}`

	env, err := NewParser().Parse(strings.NewReader(doc))
	require.NoError(t, err)

	orLater, ok := env.Graph[0].(*expandedlicensing.OrLaterOperator)
	require.True(t, ok)
	gpl, ok := orLater.SubjectLicense.(*expandedlicensing.CustomLicense)
	require.True(t, ok)
	require.Equal(t, "the GPL text", gpl.LicenseText)

	with, ok := env.Graph[1].(*expandedlicensing.WithAdditionOperator)
	require.True(t, ok)
	require.Equal(t, "https://example.com/gpl", with.SubjectExtendableLicense.GetID())
	addition, ok := with.SubjectAddition.(*expandedlicensing.CustomLicenseAddition)
	require.True(t, ok)
	require.Equal(t, "the exception text", addition.AdditionText)
}

// The model derives OrLaterOperator from ExtendableLicense, so it is valid
// wherever an extendable license is expected, and every licensing class is
// valid as a set member.
func TestLicenseHierarchy(t *testing.T) {
	var _ expandedlicensing.ExtendableLicenseDescendant = &expandedlicensing.OrLaterOperator{}
	var _ expandedlicensing.ExtendableLicenseDescendant = &expandedlicensing.CustomLicense{}
	var _ expandedlicensing.LicenseDescendant = &expandedlicensing.ListedLicense{}
	var _ expandedlicensing.LicenseAdditionDescendant = &expandedlicensing.CustomLicenseAddition{}
	var _ simplelicensing.AnyLicenseInfoDescendant = &expandedlicensing.ConjunctiveLicenseSet{}
	var _ simplelicensing.AnyLicenseInfoDescendant = &expandedlicensing.WithAdditionOperator{}
	var _ simplelicensing.AnyLicenseInfoDescendant = &simplelicensing.LicenseExpression{}
}
