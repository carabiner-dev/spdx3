package dataset

import (
	"time"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "dataset"

type Profile struct{}

func (p *Profile) Prefix() string {
	return Prefix
}

// Package dataset package extends software package with dataset fields
type Package struct {
	Node
}

func (dp *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, dp, &dp.PreNode)
}

type Node struct {
	core.Node
	BuiltTime                       *time.Time               `json:"builtTime"`
	ReleaseTime                     *time.Time               `json:"releaseTime"`
	ConfidentialityLevel            ConfidentialityLevelType `json:"confidentialityLevel"`
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
