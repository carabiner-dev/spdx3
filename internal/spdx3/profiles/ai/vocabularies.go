package ai

// EnergyUnitType vocabulary defines valid energy unit types
type EnergyUnitType string

const (
	// EnergyUnitTypeKilowattHour represents kilowatt-hour energy measurement
	EnergyUnitTypeKilowattHour EnergyUnitType = "kilowattHour"
	// EnergyUnitTypeMegajoule represents megajoule energy measurement
	EnergyUnitTypeMegajoule EnergyUnitType = "megajoule"
	// EnergyUnitTypeOther represents other energy measurement units
	EnergyUnitTypeOther EnergyUnitType = "other"
)

// SafetyRiskAssessmentType vocabulary defines valid safety risk assessment levels
type SafetyRiskAssessmentType string

const (
	// SafetyRiskAssessmentTypeHigh indicates high safety risk
	SafetyRiskAssessmentTypeHigh SafetyRiskAssessmentType = "high"
	// SafetyRiskAssessmentTypeLow indicates low safety risk
	SafetyRiskAssessmentTypeLow SafetyRiskAssessmentType = "low"
	// SafetyRiskAssessmentTypeMedium indicates medium safety risk
	SafetyRiskAssessmentTypeMedium SafetyRiskAssessmentType = "medium"
	// SafetyRiskAssessmentTypeSerious indicates serious safety risk
	SafetyRiskAssessmentTypeSerious SafetyRiskAssessmentType = "serious"
)
