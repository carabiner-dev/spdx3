package software

import (
	"time"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "software"

type Node struct {
	core.Node
	DownloadLocation string `json:"software_downloadLocation"`
	PrimaryPurpose   string `json:"software_primaryPurpose"`
}

// SoftwareArtifact represents a distinct article or unit related to software
type SoftwareArtifact struct {
	core.Artifact
	AdditionalPurpose  []string `json:"additionalPurpose,omitempty"`
	AttributionText    []string `json:"attributionText,omitempty"`
	ContentIdentifier  []string `json:"contentIdentifier,omitempty"`
	CopyrightText      string   `json:"copyrightText,omitempty"`
	Extension          []string `json:"extension,omitempty"`
	ExternalIdentifier []string `json:"externalIdentifier,omitempty"`
	ExternalRef        []string `json:"externalRef,omitempty"`
	PrimaryPurpose     string   `json:"primaryPurpose,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	VerifiedUsing      []string `json:"verifiedUsing,omitempty"`
}

func (sa *SoftwareArtifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sa, &sa.PreNode)
}

// File
type File struct {
	Node
	AdditionalPurpose  []string      `json:"additionalPurpose,omitempty"`
	AttributionText    []string      `json:"attributionText,omitempty"`
	BuiltTime          *time.Time    `json:"builtTime,omitempty"`
	ContentIdentifier  []string      `json:"contentIdentifier,omitempty"`
	ContentType        string        `json:"contentType,omitempty"`
	CopyrightText      string        `json:"copyrightText,omitempty"`
	Extension          []string      `json:"extension,omitempty"`
	ExternalIdentifier []string      `json:"externalIdentifier,omitempty"`
	ExternalRef        []string      `json:"externalRef,omitempty"`
	FileKind           string        `json:"fileKind,omitempty"`
	OriginatedBy       []*core.Agent `json:"originatedBy,omitempty"`
	ReleaseTime        *time.Time    `json:"releaseTime,omitempty"`
	StandardName       []string      `json:"standardName,omitempty"`
	Summary            string        `json:"summary,omitempty"`
	SuppliedBy         *core.Agent   `json:"suppliedBy,omitempty"`
	SupportLevel       []string      `json:"supportLevel,omitempty"`
	ValidUntilTime     *time.Time    `json:"validUntilTime,omitempty"`
	VerifiedUsing      []string      `json:"verifiedUsing,omitempty"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, f, &f.PreNode)
}

type IntegrityMethod interface {
	GetComment() string
}

// Package
type Package struct {
	Node
	AdditionalPurpose  []string          `json:"additionalPurpose,omitempty"`
	AttributionText    []string          `json:"attributionText,omitempty"`
	BuiltTime          *time.Time        `json:"builtTime,omitempty"`
	ContentIdentifier  []string          `json:"contentIdentifier,omitempty"`
	CopyrightText      string            `json:"copyrightText,omitempty"`
	Extension          []string          `json:"extension,omitempty"`
	ExternalIdentifier []string          `json:"externalIdentifier,omitempty"`
	ExternalRef        []string          `json:"externalRef,omitempty"`
	HomePage           string            `json:"homePage,omitempty"`
	OriginatedBy       []types.Node      `json:"originatedBy,omitempty"`
	PackageUrl         string            `json:"packageUrl,omitempty"`
	PackageVersion     string            `json:"packageVersion,omitempty"`
	ReleaseTime        *time.Time        `json:"releaseTime,omitempty"`
	SourceInfo         string            `json:"sourceInfo,omitempty"`
	StandardName       []string          `json:"standardName,omitempty"`
	Summary            string            `json:"summary,omitempty"`
	SuppliedBy         *core.Agent       `json:"suppliedBy,omitempty"`
	SupportLevel       []string          `json:"supportLevel,omitempty"`
	ValidUntilTime     *time.Time        `json:"validUntilTime,omitempty"`
	VerifiedUsing      []IntegrityMethod `json:"verifiedUsing,omitempty"`
}

func (p *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, p, &p.PreNode)
}
