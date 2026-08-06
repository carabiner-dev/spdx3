// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package software

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix                 = "software"
	SoftwareArtifactType   = "SoftwareArtifact"
	FileType               = "File"
	PackageType            = "Package"
	SnippetType            = "Snippet"
	SbomClass              = "Sbom"
	ContentIdentifierClass = "ContentIdentifier"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		SoftwareArtifactType:   &SoftwareArtifact{},
		FileType:               &File{},
		PackageType:            &Package{},
		SnippetType:            &Snippet{},
		SbomClass:              &Sbom{},
		ContentIdentifierClass: &ContentIdentifier{},
		fmt.Sprintf("%s_%s", Prefix, SoftwareArtifactType):   &SoftwareArtifact{},
		fmt.Sprintf("%s_%s", Prefix, FileType):               &File{},
		fmt.Sprintf("%s_%s", Prefix, PackageType):            &Package{},
		fmt.Sprintf("%s_%s", Prefix, SnippetType):            &Snippet{},
		fmt.Sprintf("%s_%s", Prefix, SbomClass):              &Sbom{},
		fmt.Sprintf("%s_%s", Prefix, ContentIdentifierClass): &ContentIdentifier{},
	},
}

// SoftwareArtifact represents a distinct article or unit related to software
type SoftwareArtifact struct {
	core.Artifact
	AdditionalPurpose []SoftwarePurpose   `json:"software_additionalPurpose,omitempty"`
	AttributionText   []string            `json:"software_attributionText,omitempty"`
	ContentIdentifier []ContentIdentifier `json:"software_contentIdentifier,omitempty"`
	CopyrightText     string              `json:"software_copyrightText,omitempty"`
	PrimaryPurpose    SoftwarePurpose     `json:"software_primaryPurpose,omitempty"`
}

func (sa *SoftwareArtifact) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, SoftwareArtifactType)
}

func (sa *SoftwareArtifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sa, &sa.PreNode)
}

// File
type File struct {
	SoftwareArtifact
	ContentType string       `json:"contentType,omitempty"`
	FileKind    FileKindType `json:"software_fileKind,omitempty"`
}

func (f *File) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, FileType)
}

func (f *File) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, f, &f.PreNode)
}

type IntegrityMethod interface {
	GetComment() string
}

// Package
type Package struct {
	SoftwareArtifact
	DownloadLocation string `json:"software_downloadLocation,omitempty"`
	HomePage         string `json:"software_homePage,omitempty"`
	PackageUrl       string `json:"software_packageUrl,omitempty"`
	PackageVersion   string `json:"software_packageVersion,omitempty"`
	SourceInfo       string `json:"software_sourceInfo,omitempty"`
}

func (p *Package) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, PackageType)
}

func (p *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, p, &p.PreNode)
}

// ContentIdentifier provides a canonical, unique, immutable identifier
type ContentIdentifier struct {
	core.IntegrityMethod
	ContentIdentifierType  ContentIdentifierType `json:"software_contentIdentifierType"`
	ContentIdentifierValue string                `json:"software_contentIdentifierValue"`
}

func (ci *ContentIdentifier) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, ContentIdentifierClass)
}

// Snippet represents a portion of a file
type Snippet struct {
	SoftwareArtifact
	ByteRange       *core.PositiveIntegerRange `json:"software_byteRange,omitempty"`
	LineRange       *core.PositiveIntegerRange `json:"software_lineRange,omitempty"`
	SnippetFromFile string                     `json:"software_snippetFromFile"`
}

func (s *Snippet) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, SnippetType)
}

func (s *Snippet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}

// Sbom represents a Software Bill of Materials
type Sbom struct {
	core.Bom
	SbomType []SbomType `json:"software_sbomType,omitempty"`
}

func (s *Sbom) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, SbomClass)
}

func (s *Sbom) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}
