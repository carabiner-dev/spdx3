package software

import (
	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "software"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"SoftwareArtifact":          &SoftwareArtifact{},
		"File":                      &File{},
		"Package":                   &Package{},
		"Snippet":                   &Snippet{},
		"Sbom":                      &Sbom{},
		"software_SoftwareArtifact": &SoftwareArtifact{},
		"software_File":             &File{},
		"software_Package":          &Package{},
		"software_Snippet":          &Snippet{},
		"software_Sbom":             &Sbom{},
	},
}

// SoftwareArtifact represents a distinct article or unit related to software
type SoftwareArtifact struct {
	core.Artifact
	AdditionalPurpose  []SoftwarePurpose `json:"additionalPurpose,omitempty"`
	AttributionText    []string          `json:"attributionText,omitempty"`
	ContentIdentifier  []string          `json:"contentIdentifier,omitempty"`
	CopyrightText      string            `json:"copyrightText,omitempty"`
	Extension          []string          `json:"extension,omitempty"`
	ExternalIdentifier []string          `json:"externalIdentifier,omitempty"`
	ExternalRef        []string          `json:"externalRef,omitempty"`
	PrimaryPurpose     SoftwarePurpose   `json:"primaryPurpose,omitempty"`
	Summary            string            `json:"summary,omitempty"`
	VerifiedUsing      []string          `json:"verifiedUsing,omitempty"`
}

func (sa *SoftwareArtifact) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, sa, &sa.PreNode)
}

// File
type File struct {
	SoftwareArtifact
	ContentType string       `json:"contentType,omitempty"`
	FileKind    FileKindType `json:"fileKind,omitempty"`
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
	DownloadLocation string `json:"downloadLocation,omitempty"`
	HomePage         string `json:"homePage,omitempty"`
	PackageUrl       string `json:"packageUrl,omitempty"`
	PackageVersion   string `json:"packageVersion,omitempty"`
	SourceInfo       string `json:"sourceInfo,omitempty"`
}

func (p *Package) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, p, &p.PreNode)
}

// ContentIdentifier provides a canonical, unique, immutable identifier
type ContentIdentifier struct {
	core.IntegrityMethod
	ContentIdentifierType  ContentIdentifierType `json:"contentIdentifierType"`
	ContentIdentifierValue string                `json:"contentIdentifierValue"`
}

// Snippet represents a portion of a file
type Snippet struct {
	SoftwareArtifact
	ByteRange       *core.PositiveIntegerRange `json:"byteRange,omitempty"`
	LineRange       *core.PositiveIntegerRange `json:"lineRange,omitempty"`
	SnippetFromFile string                     `json:"snippetFromFile"`
}

func (s *Snippet) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}

// Sbom represents a Software Bill of Materials
type Sbom struct {
	core.Bom
	SbomType []SbomType `json:"sbomType,omitempty"`
}

func (s *Sbom) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, s, &s.PreNode)
}
