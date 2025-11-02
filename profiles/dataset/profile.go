// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"time"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const Prefix = "dataset"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"DatasetPackage":         &Package{},
		"dataset_DatasetPackage": &Package{},
	},
}

// Package dataset package extends software package with dataset fields
type Package struct {
	core.Element
	BuiltTime                       *time.Time               `json:"builtTime"`
	ReleaseTime                     *time.Time               `json:"releaseTime"`
	ConfidentialityLevel            ConfidentialityLevelType `json:"dataset_confidentialityLevel"`
	DataPreprocessing               []string                 `json:"dataset_dataPreprocessing"`
	DatasetAvailability             DatasetAvailabilityType  `json:"dataset_datasetAvailability"`
	DataCollectionProcess           string                   `json:"dataset_dataCollectionProcess"`
	DatasetSize                     int64                    `json:"dataset_datasetSize"`
	DatasetType                     []DatasetType            `json:"dataset_datasetType"`
	DatasetUpdateMechanism          string                   `json:"dataset_datasetUpdateMechanism"`
	HasSensitivePersonalInformation string                   `json:"dataset_hasSensitivePersonalInformation"`
	IntendedUse                     string                   `json:"dataset_intendedUse"`
	KnownBias                       []string                 `json:"dataset_knownBias"`
}

func (dp *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, dp, &dp.PreNode)
}
