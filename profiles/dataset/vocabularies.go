// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package dataset

import "github.com/carabiner-dev/spdx3/types"

// ConfidentialityLevelType vocabulary defines valid confidentiality levels using Traffic Light Protocol
type ConfidentialityLevelType string

var ConfidentialityLevelTypes = types.Vocabulary[ConfidentialityLevelType]{
	ConfidentialityLevelTypeAmber,
	ConfidentialityLevelTypeClear,
	ConfidentialityLevelTypeGreen,
	ConfidentialityLevelTypeRed,
}

const (
	// ConfidentialityLevelTypeAmber indicates limited disclosure to participants' organizations
	ConfidentialityLevelTypeAmber ConfidentialityLevelType = "amber"
	// ConfidentialityLevelTypeClear indicates unlimited disclosure
	ConfidentialityLevelTypeClear ConfidentialityLevelType = "clear"
	// ConfidentialityLevelTypeGreen indicates limited disclosure within community
	ConfidentialityLevelTypeGreen ConfidentialityLevelType = "green"
	// ConfidentialityLevelTypeRed indicates no disclosure beyond participants
	ConfidentialityLevelTypeRed ConfidentialityLevelType = "red"
)

// DatasetAvailabilityType vocabulary defines valid dataset availability types
type DatasetAvailabilityType string

var DatasetAvailabilityTypes = types.Vocabulary[DatasetAvailabilityType]{
	DatasetAvailabilityTypeClickthrough,
	DatasetAvailabilityTypeDirectDownload,
	DatasetAvailabilityTypeQuery,
	DatasetAvailabilityTypeRegistration,
	DatasetAvailabilityTypeScrapingScript,
}

const (
	// DatasetAvailabilityTypeClickthrough indicates dataset requires accepting terms through web interface
	DatasetAvailabilityTypeClickthrough DatasetAvailabilityType = "clickthrough"
	// DatasetAvailabilityTypeDirectDownload indicates dataset is directly downloadable without restrictions
	DatasetAvailabilityTypeDirectDownload DatasetAvailabilityType = "directDownload"
	// DatasetAvailabilityTypeQuery indicates dataset requires query to access portions
	DatasetAvailabilityTypeQuery DatasetAvailabilityType = "query"
	// DatasetAvailabilityTypeRegistration indicates dataset requires registration before access
	DatasetAvailabilityTypeRegistration DatasetAvailabilityType = "registration"
	// DatasetAvailabilityTypeScrapingScript indicates dataset requires script to collect data
	DatasetAvailabilityTypeScrapingScript DatasetAvailabilityType = "scrapingScript"
)

// DatasetType vocabulary defines valid dataset types
type DatasetType string

var DatasetTypes = types.Vocabulary[DatasetType]{
	DatasetTypeAudio,
	DatasetTypeCategorical,
	DatasetTypeGraph,
	DatasetTypeImage,
	DatasetTypeNoAssertion,
	DatasetTypeNumeric,
	DatasetTypeOther,
	DatasetTypeSensor,
	DatasetTypeStructured,
	DatasetTypeSyntactic,
	DatasetTypeText,
	DatasetTypeTimeseries,
	DatasetTypeTimestamp,
	DatasetTypeVideo,
}

const (
	// DatasetTypeAudio indicates dataset contains audio data
	DatasetTypeAudio DatasetType = "audio"
	// DatasetTypeCategorical indicates data classified into a discrete number of categories
	DatasetTypeCategorical DatasetType = "categorical"
	// DatasetTypeGraph indicates dataset is in graph format
	DatasetTypeGraph DatasetType = "graph"
	// DatasetTypeImage indicates dataset contains image data
	DatasetTypeImage DatasetType = "image"
	// DatasetTypeNoAssertion makes no assertion about dataset type
	DatasetTypeNoAssertion DatasetType = "noAssertion"
	// DatasetTypeNumeric indicates dataset contains numeric data
	DatasetTypeNumeric DatasetType = "numeric"
	// DatasetTypeOther indicates dataset is of another type
	DatasetTypeOther DatasetType = "other"
	// DatasetTypeSensor indicates dataset contains sensor data
	DatasetTypeSensor DatasetType = "sensor"
	// DatasetTypeStructured indicates dataset has structured format
	DatasetTypeStructured DatasetType = "structured"
	// DatasetTypeSyntactic indicates data describing the syntax or semantics of a language
	DatasetTypeSyntactic DatasetType = "syntactic"
	// DatasetTypeText indicates dataset contains text data
	DatasetTypeText DatasetType = "text"
	// DatasetTypeTimeseries indicates data recorded as an ordered sequence of timestamped entries
	DatasetTypeTimeseries DatasetType = "timeseries"
	// DatasetTypeTimestamp indicates data recorded with a timestamp per entry, not necessarily ordered
	DatasetTypeTimestamp DatasetType = "timestamp"
	// DatasetTypeVideo indicates dataset contains video data
	DatasetTypeVideo DatasetType = "video"
)

// IsValid reports whether the value is a member of its vocabulary, so that
// the parser and Validate can check it without knowing which one it is.
func (c ConfidentialityLevelType) IsValid() bool { return ConfidentialityLevelTypes.Contains(c) }
func (d DatasetAvailabilityType) IsValid() bool  { return DatasetAvailabilityTypes.Contains(d) }
func (d DatasetType) IsValid() bool              { return DatasetTypes.Contains(d) }
