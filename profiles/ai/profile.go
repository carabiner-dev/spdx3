// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package ai

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix        = "ai"
	AIPackageType = "AIPackage"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		AIPackageType: &AIPackage{},
		fmt.Sprintf("%s_%s", Prefix, AIPackageType): &AIPackage{},
	},
}

// EnergyConsumptionDescription provides detailed documentation of energy consumption metrics
type EnergyConsumptionDescription struct {
	EnergyQuantity types.Decimal  `json:"ai_energyQuantity"`
	EnergyUnit     EnergyUnitType `json:"ai_energyUnit"`
}

// EnergyConsumption captures information about energy usage associated with AI model operations
type EnergyConsumption struct {
	TrainingEnergyConsumption   []EnergyConsumptionDescription `json:"ai_trainingEnergyConsumption,omitempty"`
	FinetuningEnergyConsumption []EnergyConsumptionDescription `json:"ai_finetuningEnergyConsumption,omitempty"`
	InferenceEnergyConsumption  []EnergyConsumptionDescription `json:"ai_inferenceEnergyConsumption,omitempty"`
}

// AIPackage represents an AI software package or system
type AIPackage struct {
	software.Package
	AutonomyType                    core.PresenceType        `json:"ai_autonomyType,omitempty"`
	Domain                          []string                 `json:"ai_domain,omitempty"`
	EnergyConsumption               *EnergyConsumption       `json:"ai_energyConsumption,omitempty"`
	Hyperparameter                  []core.DictionaryEntry   `json:"ai_hyperparameter,omitempty"`
	InformationAboutApplication     string                   `json:"ai_informationAboutApplication,omitempty"`
	InformationAboutTraining        string                   `json:"ai_informationAboutTraining,omitempty"`
	Limitation                      string                   `json:"ai_limitation,omitempty"`
	Metric                          []core.DictionaryEntry   `json:"ai_metric,omitempty"`
	MetricDecisionThreshold         []core.DictionaryEntry   `json:"ai_metricDecisionThreshold,omitempty"`
	ModelDataPreprocessing          []string                 `json:"ai_modelDataPreprocessing,omitempty"`
	ModelExplainability             []string                 `json:"ai_modelExplainability,omitempty"`
	SafetyRiskAssessment            SafetyRiskAssessmentType `json:"ai_safetyRiskAssessment,omitempty"`
	StandardCompliance              []string                 `json:"ai_standardCompliance,omitempty"`
	TypeOfModel                     []string                 `json:"ai_typeOfModel,omitempty"`
	UseSensitivePersonalInformation core.PresenceType        `json:"ai_useSensitivePersonalInformation,omitempty"`
}

func (a *AIPackage) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, AIPackageType)
}

func (a *AIPackage) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}
