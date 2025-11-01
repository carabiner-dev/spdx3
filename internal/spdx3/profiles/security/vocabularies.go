package security

// CvssSeverityType vocabulary defines valid CVSS severity levels
type CvssSeverityType string

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

const (
	// ExploitCatalogTypeKev represents CISA's Known Exploited Vulnerabilities catalog
	ExploitCatalogTypeKev ExploitCatalogType = "kev"
	// ExploitCatalogTypeOther represents other exploit catalogs
	ExploitCatalogTypeOther ExploitCatalogType = "other"
)

// SsvcDecisionType vocabulary defines valid SSVC (Stakeholder-Specific Vulnerability Categorization) decision types
type SsvcDecisionType string

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

const (
	// VexJustificationTypeComponentNotPresent indicates the vulnerable component is not present
	VexJustificationTypeComponentNotPresent VexJustificationType = "componentNotPresent"
	// VexJustificationTypeInlineMitigationsAlreadyExist indicates built-in protections prevent exploitation
	VexJustificationTypeInlineMitigationsAlreadyExist VexJustificationType = "inlineMitigationsAlreadyExist"
	// VexJustificationTypeRequiresConfiguration indicates vulnerability requires specific configuration
	VexJustificationTypeRequiresConfiguration VexJustificationType = "requiresConfiguration"
	// VexJustificationTypeRequiresDependency indicates vulnerability requires specific dependency
	VexJustificationTypeRequiresDependency VexJustificationType = "requiresDependency"
	// VexJustificationTypeVulnerableCodeCannotBeControlledByAdversary indicates adversary cannot control vulnerable code
	VexJustificationTypeVulnerableCodeCannotBeControlledByAdversary VexJustificationType = "vulnerableCodeCannotBeControlledByAdversary"
	// VexJustificationTypeVulnerableCodeNotInExecutePath indicates vulnerable code is not in execution path
	VexJustificationTypeVulnerableCodeNotInExecutePath VexJustificationType = "vulnerableCodeNotInExecutePath"
	// VexJustificationTypeVulnerableCodeNotPresent indicates vulnerable code is not present in component
	VexJustificationTypeVulnerableCodeNotPresent VexJustificationType = "vulnerableCodeNotPresent"
)
