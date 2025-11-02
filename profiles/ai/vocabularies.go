package ai

import "github.com/carabiner-dev/spdx3/types"

// EnergyUnitType vocabulary defines valid energy unit types
type EnergyUnitType string

var EnergyUnitTypes = types.Vocabulary[EnergyUnitType]{
	EnergyUnitTypeKilowattHour,
	EnergyUnitTypeMegajoule,
	EnergyUnitTypeOther,
}

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

var SafetyRiskAssessmentTypes = types.Vocabulary[SafetyRiskAssessmentType]{
	SafetyRiskAssessmentTypeHigh,
	SafetyRiskAssessmentTypeLow,
	SafetyRiskAssessmentTypeMedium,
	SafetyRiskAssessmentTypeSerious,
}

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
