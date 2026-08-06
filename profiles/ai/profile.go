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
		AIPackageType:                          &AIPackage{},
		fmt.Sprintf("%s_%s", Prefix, AIPackageType): &AIPackage{},
	},
}

// EnergyConsumptionDescription provides detailed documentation of energy consumption metrics
type EnergyConsumptionDescription struct {
	EnergyQuantity float64        `json:"energyQuantity"`
	EnergyUnit     EnergyUnitType `json:"energyUnit"`
}

// EnergyConsumption captures information about energy usage associated with AI model operations
type EnergyConsumption struct {
	TrainingEnergyConsumption   []EnergyConsumptionDescription `json:"trainingEnergyConsumption,omitempty"`
	FinetuningEnergyConsumption []EnergyConsumptionDescription `json:"finetuningEnergyConsumption,omitempty"`
	InferenceEnergyConsumption  []EnergyConsumptionDescription `json:"inferenceEnergyConsumption,omitempty"`
}

// AIPackage represents an AI software package or system
type AIPackage struct {
	software.Package
	AutonomyType                    core.PresenceType        `json:"autonomyType,omitempty"`
	Domain                          []string                 `json:"domain,omitempty"`
	EnergyConsumption               *EnergyConsumption       `json:"energyConsumption,omitempty"`
	Hyperparameter                  []core.DictionaryEntry   `json:"hyperparameter,omitempty"`
	InformationAboutApplication     string                   `json:"informationAboutApplication,omitempty"`
	InformationAboutTraining        string                   `json:"informationAboutTraining,omitempty"`
	Limitation                      string                   `json:"limitation,omitempty"`
	Metric                          []core.DictionaryEntry   `json:"metric,omitempty"`
	MetricDecisionThreshold         []core.DictionaryEntry   `json:"metricDecisionThreshold,omitempty"`
	ModelDataPreprocessing          []string                 `json:"modelDataPreprocessing,omitempty"`
	ModelExplainability             []string                 `json:"modelExplainability,omitempty"`
	SafetyRiskAssessment            SafetyRiskAssessmentType `json:"safetyRiskAssessment,omitempty"`
	StandardCompliance              []string                 `json:"standardCompliance,omitempty"`
	TypeOfModel                     []string                 `json:"typeOfModel,omitempty"`
	UseSensitivePersonalInformation core.PresenceType        `json:"useSensitivePersonalInformation,omitempty"`
}

func (a *AIPackage) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, AIPackageType)
}

func (a *AIPackage) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, a, &a.PreNode)
}
