// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix             = "dataset"
	DatasetPackageType = "DatasetPackage"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		DatasetPackageType: &Package{},
		fmt.Sprintf("%s_%s", Prefix, DatasetPackageType): &Package{},
	},
}

// Package dataset package extends software package with dataset fields
type Package struct {
	software.Package
	AnonymizationMethodUsed         []string                 `json:"dataset_anonymizationMethodUsed,omitempty"`
	ConfidentialityLevel            ConfidentialityLevelType `json:"dataset_confidentialityLevel,omitempty"`
	DataCollectionProcess           string                   `json:"dataset_dataCollectionProcess,omitempty"`
	DataPreprocessing               []string                 `json:"dataset_dataPreprocessing,omitempty"`
	DatasetAvailability             DatasetAvailabilityType  `json:"dataset_datasetAvailability,omitempty"`
	DatasetNoise                    string                   `json:"dataset_datasetNoise,omitempty"`
	DatasetSize                     int64                    `json:"dataset_datasetSize,omitempty"`
	DatasetType                     []DatasetType            `json:"dataset_datasetType,omitempty"`
	DatasetUpdateMechanism          string                   `json:"dataset_datasetUpdateMechanism,omitempty"`
	HasSensitivePersonalInformation string                   `json:"dataset_hasSensitivePersonalInformation,omitempty"`
	IntendedUse                     string                   `json:"dataset_intendedUse,omitempty"`
	KnownBias                       []string                 `json:"dataset_knownBias,omitempty"`
	Sensor                          []core.DictionaryEntry   `json:"dataset_sensor,omitempty"`
}

func (dp *Package) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, DatasetPackageType)
}

func (dp *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, dp, &dp.PreNode)
}
