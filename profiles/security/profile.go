// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package security

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                                              = "security"
	VulnerabilityType                                   = "Vulnerability"
	CvssV2VulnAssessmentRelationshipType                = "CvssV2VulnAssessmentRelationship"
	CvssV3VulnAssessmentRelationshipType                = "CvssV3VulnAssessmentRelationship"
	CvssV4VulnAssessmentRelationshipType                = "CvssV4VulnAssessmentRelationship"
	EpssVulnAssessmentRelationshipType                  = "EpssVulnAssessmentRelationship"
	SsvcVulnAssessmentRelationshipType                  = "SsvcVulnAssessmentRelationship"
	ExploitCatalogVulnAssessmentRelationshipType        = "ExploitCatalogVulnAssessmentRelationship"
	VexAffectedVulnAssessmentRelationshipType           = "VexAffectedVulnAssessmentRelationship"
	VexFixedVulnAssessmentRelationshipType              = "VexFixedVulnAssessmentRelationship"
	VexNotAffectedVulnAssessmentRelationshipType        = "VexNotAffectedVulnAssessmentRelationship"
	VexUnderInvestigationVulnAssessmentRelationshipType = "VexUnderInvestigationVulnAssessmentRelationship"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		VulnerabilityType:                                                                 &Vulnerability{},
		CvssV2VulnAssessmentRelationshipType:                                              &CvssV2VulnAssessmentRelationship{},
		CvssV3VulnAssessmentRelationshipType:                                              &CvssV3VulnAssessmentRelationship{},
		CvssV4VulnAssessmentRelationshipType:                                              &CvssV4VulnAssessmentRelationship{},
		EpssVulnAssessmentRelationshipType:                                                &EpssVulnAssessmentRelationship{},
		SsvcVulnAssessmentRelationshipType:                                                &SsvcVulnAssessmentRelationship{},
		ExploitCatalogVulnAssessmentRelationshipType:                                      &ExploitCatalogVulnAssessmentRelationship{},
		VexAffectedVulnAssessmentRelationshipType:                                         &VexAffectedVulnAssessmentRelationship{},
		VexFixedVulnAssessmentRelationshipType:                                            &VexFixedVulnAssessmentRelationship{},
		VexNotAffectedVulnAssessmentRelationshipType:                                      &VexNotAffectedVulnAssessmentRelationship{},
		VexUnderInvestigationVulnAssessmentRelationshipType:                               &VexUnderInvestigationVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, VulnerabilityType):                                   &Vulnerability{},
		fmt.Sprintf("%s_%s", Prefix, CvssV2VulnAssessmentRelationshipType):                &CvssV2VulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, CvssV3VulnAssessmentRelationshipType):                &CvssV3VulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, CvssV4VulnAssessmentRelationshipType):                &CvssV4VulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, EpssVulnAssessmentRelationshipType):                  &EpssVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, SsvcVulnAssessmentRelationshipType):                  &SsvcVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, ExploitCatalogVulnAssessmentRelationshipType):        &ExploitCatalogVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, VexAffectedVulnAssessmentRelationshipType):           &VexAffectedVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, VexFixedVulnAssessmentRelationshipType):              &VexFixedVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, VexNotAffectedVulnAssessmentRelationshipType):        &VexNotAffectedVulnAssessmentRelationship{},
		fmt.Sprintf("%s_%s", Prefix, VexUnderInvestigationVulnAssessmentRelationshipType): &VexUnderInvestigationVulnAssessmentRelationship{},
	},
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	core.Artifact
	ModifiedTime  *types.DateTime `json:"security_modifiedTime,omitempty"`
	PublishedTime *types.DateTime `json:"security_publishedTime,omitempty"`
	WithdrawnTime *types.DateTime `json:"security_withdrawnTime,omitempty"`
}

func (v *Vulnerability) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, VulnerabilityType)
}

func (v *Vulnerability) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VulnAssessmentRelationship is the abstract ancestor for all vulnerability assessments
type VulnAssessmentRelationship struct {
	core.Relationship
	AssessedElement string               `json:"security_assessedElement,omitempty"`
	ModifiedTime    *types.DateTime      `json:"security_modifiedTime,omitempty"`
	PublishedTime   *types.DateTime      `json:"security_publishedTime,omitempty"`
	SuppliedBy      core.AgentDescendant `json:"suppliedBy,omitempty"`
	WithdrawnTime   *types.DateTime      `json:"security_withdrawnTime,omitempty"`
}

// CvssV2VulnAssessmentRelationship represents a CVSS v2 assessment
type CvssV2VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        types.Decimal `json:"security_score"`
	VectorString string        `json:"security_vectorString"`
}

func (c *CvssV2VulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CvssV2VulnAssessmentRelationshipType)
}

func (c *CvssV2VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// CvssV3VulnAssessmentRelationship represents a CVSS v3 assessment
type CvssV3VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        types.Decimal    `json:"security_score"`
	Severity     CvssSeverityType `json:"security_severity"`
	VectorString string           `json:"security_vectorString"`
}

func (c *CvssV3VulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CvssV3VulnAssessmentRelationshipType)
}

func (c *CvssV3VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// CvssV4VulnAssessmentRelationship represents a CVSS v4 assessment
type CvssV4VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        types.Decimal    `json:"security_score"`
	Severity     CvssSeverityType `json:"security_severity"`
	VectorString string           `json:"security_vectorString"`
}

func (c *CvssV4VulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, CvssV4VulnAssessmentRelationshipType)
}

func (c *CvssV4VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// EpssVulnAssessmentRelationship represents an EPSS assessment
type EpssVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Percentile  types.Decimal `json:"security_percentile"`
	Probability types.Decimal `json:"security_probability"`
}

func (e *EpssVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, EpssVulnAssessmentRelationshipType)
}

func (e *EpssVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}

// SsvcVulnAssessmentRelationship represents an SSVC assessment
type SsvcVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	DecisionType SsvcDecisionType `json:"security_decisionType"`
}

func (s *SsvcVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, SsvcVulnAssessmentRelationshipType)
}

func (s *SsvcVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}

// ExploitCatalogVulnAssessmentRelationship references known exploit catalogs
type ExploitCatalogVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	CatalogType ExploitCatalogType `json:"security_catalogType"`
	Exploited   bool               `json:"security_exploited"`
	Locator     string             `json:"security_locator"`
}

func (e *ExploitCatalogVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, ExploitCatalogVulnAssessmentRelationshipType)
}

func (e *ExploitCatalogVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}

// VexVulnAssessmentRelationship is the abstract base for VEX assessments
type VexVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	StatusNotes string `json:"security_statusNotes,omitempty"`
	VexVersion  string `json:"security_vexVersion,omitempty"`
}

// VexAffectedVulnAssessmentRelationship represents a VEX affected statement
type VexAffectedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
	ActionStatement     string          `json:"security_actionStatement"`
	ActionStatementTime *types.DateTime `json:"security_actionStatementTime,omitempty"`
}

func (v *VexAffectedVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, VexAffectedVulnAssessmentRelationshipType)
}

func (v *VexAffectedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexFixedVulnAssessmentRelationship represents a VEX fixed statement
type VexFixedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
}

func (v *VexFixedVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, VexFixedVulnAssessmentRelationshipType)
}

func (v *VexFixedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexNotAffectedVulnAssessmentRelationship represents a VEX not_affected statement
type VexNotAffectedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
	ImpactStatement     string               `json:"security_impactStatement,omitempty"`
	ImpactStatementTime *types.DateTime      `json:"security_impactStatementTime,omitempty"`
	JustificationType   VexJustificationType `json:"security_justificationType,omitempty"`
}

func (v *VexNotAffectedVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, VexNotAffectedVulnAssessmentRelationshipType)
}

func (v *VexNotAffectedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexUnderInvestigationVulnAssessmentRelationship represents a VEX under_investigation statement
type VexUnderInvestigationVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
}

func (v *VexUnderInvestigationVulnAssessmentRelationship) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, VexUnderInvestigationVulnAssessmentRelationshipType)
}

func (v *VexUnderInvestigationVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}
