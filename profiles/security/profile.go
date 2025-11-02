package security

import (
	"time"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const Prefix = "security"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"Vulnerability":                                   &Vulnerability{},
		"CvssV2VulnAssessmentRelationship":                &CvssV2VulnAssessmentRelationship{},
		"CvssV3VulnAssessmentRelationship":                &CvssV3VulnAssessmentRelationship{},
		"CvssV4VulnAssessmentRelationship":                &CvssV4VulnAssessmentRelationship{},
		"EpssVulnAssessmentRelationship":                  &EpssVulnAssessmentRelationship{},
		"SsvcVulnAssessmentRelationship":                  &SsvcVulnAssessmentRelationship{},
		"ExploitCatalogVulnAssessmentRelationship":        &ExploitCatalogVulnAssessmentRelationship{},
		"VexAffectedVulnAssessmentRelationship":           &VexAffectedVulnAssessmentRelationship{},
		"VexFixedVulnAssessmentRelationship":              &VexFixedVulnAssessmentRelationship{},
		"VexNotAffectedVulnAssessmentRelationship":        &VexNotAffectedVulnAssessmentRelationship{},
		"VexUnderInvestigationVulnAssessmentRelationship": &VexUnderInvestigationVulnAssessmentRelationship{},
	},
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	core.Artifact
	ModifiedTime  *time.Time `json:"modifiedTime,omitempty"`
	PublishedTime *time.Time `json:"publishedTime,omitempty"`
	WithdrawnTime *time.Time `json:"withdrawnTime,omitempty"`
}

func (v *Vulnerability) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VulnAssessmentRelationship is the abstract ancestor for all vulnerability assessments
type VulnAssessmentRelationship struct {
	core.Relationship
	AssessedElement string      `json:"assessedElement,omitempty"`
	ModifiedTime    *time.Time  `json:"modifiedTime,omitempty"`
	PublishedTime   *time.Time  `json:"publishedTime,omitempty"`
	SuppliedBy      *core.Agent `json:"suppliedBy,omitempty"`
	WithdrawnTime   *time.Time  `json:"withdrawnTime,omitempty"`
}

// CvssV2VulnAssessmentRelationship represents a CVSS v2 assessment
type CvssV2VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        float64 `json:"score"`
	VectorString string  `json:"vectorString"`
}

func (c *CvssV2VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// CvssV3VulnAssessmentRelationship represents a CVSS v3 assessment
type CvssV3VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        float64          `json:"score"`
	Severity     CvssSeverityType `json:"severity"`
	VectorString string           `json:"vectorString"`
}

func (c *CvssV3VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// CvssV4VulnAssessmentRelationship represents a CVSS v4 assessment
type CvssV4VulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Score        float64          `json:"score"`
	Severity     CvssSeverityType `json:"severity"`
	VectorString string           `json:"vectorString"`
}

func (c *CvssV4VulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, c, &c.PreNode)
}

// EpssVulnAssessmentRelationship represents an EPSS assessment
type EpssVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	Percentile  float64 `json:"percentile"`
	Probability float64 `json:"probability"`
}

func (e *EpssVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}

// SsvcVulnAssessmentRelationship represents an SSVC assessment
type SsvcVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	DecisionType SsvcDecisionType `json:"decisionType"`
}

func (s *SsvcVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}

// ExploitCatalogVulnAssessmentRelationship references known exploit catalogs
type ExploitCatalogVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	CatalogType ExploitCatalogType `json:"catalogType"`
	Exploited   bool               `json:"exploited"`
	Locator     string             `json:"locator"`
}

func (e *ExploitCatalogVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, e, &e.PreNode)
}

// VexVulnAssessmentRelationship is the abstract base for VEX assessments
type VexVulnAssessmentRelationship struct {
	VulnAssessmentRelationship
	StatusNotes string `json:"statusNotes,omitempty"`
	VexVersion  string `json:"vexVersion,omitempty"`
}

// VexAffectedVulnAssessmentRelationship represents a VEX affected statement
type VexAffectedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
	ActionStatement     string     `json:"actionStatement"`
	ActionStatementTime *time.Time `json:"actionStatementTime,omitempty"`
}

func (v *VexAffectedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexFixedVulnAssessmentRelationship represents a VEX fixed statement
type VexFixedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
}

func (v *VexFixedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexNotAffectedVulnAssessmentRelationship represents a VEX not_affected statement
type VexNotAffectedVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
	ImpactStatement     string               `json:"impactStatement,omitempty"`
	ImpactStatementTime *time.Time           `json:"impactStatementTime,omitempty"`
	JustificationType   VexJustificationType `json:"justificationType,omitempty"`
}

func (v *VexNotAffectedVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}

// VexUnderInvestigationVulnAssessmentRelationship represents a VEX under_investigation statement
type VexUnderInvestigationVulnAssessmentRelationship struct {
	VexVulnAssessmentRelationship
}

func (v *VexUnderInvestigationVulnAssessmentRelationship) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, v, &v.PreNode)
}
