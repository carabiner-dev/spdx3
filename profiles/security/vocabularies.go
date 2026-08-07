// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package security

import "github.com/carabiner-dev/spdx3/types"

// CvssSeverityType vocabulary defines valid CVSS severity levels
type CvssSeverityType string

var CvssSeverityTypes = types.Vocabulary[CvssSeverityType]{
	CvssSeverityTypeCritical,
	CvssSeverityTypeHigh,
	CvssSeverityTypeMedium,
	CvssSeverityTypeLow,
	CvssSeverityTypeNone,
}

const (
	// CvssSeverityTypeCritical indicates critical severity
	CvssSeverityTypeCritical CvssSeverityType = "critical"
	// CvssSeverityTypeHigh indicates high severity
	CvssSeverityTypeHigh CvssSeverityType = "high"
	// CvssSeverityTypeMedium indicates medium severity
	CvssSeverityTypeMedium CvssSeverityType = "medium"
	// CvssSeverityTypeLow indicates low severity
	CvssSeverityTypeLow CvssSeverityType = "low"
	// CvssSeverityTypeNone indicates no severity
	CvssSeverityTypeNone CvssSeverityType = "none"
)

// ExploitCatalogType vocabulary defines valid exploit catalog types
type ExploitCatalogType string

var ExploitCatalogTypes = types.Vocabulary[ExploitCatalogType]{
	ExploitCatalogTypeKev,
	ExploitCatalogTypeOther,
}

const (
	// ExploitCatalogTypeKev represents CISA's Known Exploited Vulnerabilities catalog
	ExploitCatalogTypeKev ExploitCatalogType = "kev"
	// ExploitCatalogTypeOther represents other exploit catalogs
	ExploitCatalogTypeOther ExploitCatalogType = "other"
)

// SsvcDecisionType vocabulary defines valid SSVC (Stakeholder-Specific Vulnerability Categorization) decision types
type SsvcDecisionType string

var SsvcDecisionTypes = types.Vocabulary[SsvcDecisionType]{
	SsvcDecisionTypeAct,
	SsvcDecisionTypeAttend,
	SsvcDecisionTypeTrack,
	SsvcDecisionTypeTrackStar,
}

const (
	// SsvcDecisionTypeAct indicates immediate action required
	SsvcDecisionTypeAct SsvcDecisionType = "act"
	// SsvcDecisionTypeAttend indicates attention needed in scheduled maintenance
	SsvcDecisionTypeAttend SsvcDecisionType = "attend"
	// SsvcDecisionTypeTrack indicates continued monitoring with less urgency
	SsvcDecisionTypeTrack SsvcDecisionType = "track"
	// SsvcDecisionTypeTrackStar indicates tracking with more vigilance
	SsvcDecisionTypeTrackStar SsvcDecisionType = "trackStar"
)

// VexJustificationType vocabulary defines valid VEX justification types for not-affected status
type VexJustificationType string

var VexJustificationTypes = types.Vocabulary[VexJustificationType]{
	VexJustificationTypeComponentNotPresent,
	VexJustificationTypeInlineMitigationsAlreadyExist,
	VexJustificationTypeVulnerableCodeCannotBeControlledByAdversary,
	VexJustificationTypeVulnerableCodeNotInExecutePath,
	VexJustificationTypeVulnerableCodeNotPresent,
}

const (
	// VexJustificationTypeComponentNotPresent indicates the vulnerable component is not present
	VexJustificationTypeComponentNotPresent VexJustificationType = "componentNotPresent"
	// VexJustificationTypeInlineMitigationsAlreadyExist indicates built-in protections prevent exploitation
	VexJustificationTypeInlineMitigationsAlreadyExist VexJustificationType = "inlineMitigationsAlreadyExist"
	// VexJustificationTypeVulnerableCodeCannotBeControlledByAdversary indicates adversary cannot control vulnerable code
	VexJustificationTypeVulnerableCodeCannotBeControlledByAdversary VexJustificationType = "vulnerableCodeCannotBeControlledByAdversary"
	// VexJustificationTypeVulnerableCodeNotInExecutePath indicates vulnerable code is not in execution path
	VexJustificationTypeVulnerableCodeNotInExecutePath VexJustificationType = "vulnerableCodeNotInExecutePath"
	// VexJustificationTypeVulnerableCodeNotPresent indicates vulnerable code is not present in component
	VexJustificationTypeVulnerableCodeNotPresent VexJustificationType = "vulnerableCodeNotPresent"
)

// IsValid reports whether the value is a member of its vocabulary, so that
// the parser and Validate can check it without knowing which one it is.
func (c CvssSeverityType) IsValid() bool     { return CvssSeverityTypes.Contains(c) }
func (e ExploitCatalogType) IsValid() bool   { return ExploitCatalogTypes.Contains(e) }
func (s SsvcDecisionType) IsValid() bool     { return SsvcDecisionTypes.Contains(s) }
func (v VexJustificationType) IsValid() bool { return VexJustificationTypes.Contains(v) }
