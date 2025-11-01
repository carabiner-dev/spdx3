package software

import (
	"time"

	"github.com/carabiner-dev/databom/internal/spdx3/profiles/core"
	"github.com/carabiner-dev/databom/internal/spdx3/unmarshal"
)

const Prefix = "software"

type Node struct {
	core.Node
	DownloadLocation string `json:"software_downloadLocation"`
	PrimaryPurpose   string `json:"software_primaryPurpose"`
}

// File
type File struct {
	Node
	AdditionalPurpose  []string     `json:"additionalPurpose,omitempty"`
	AttributionText    []string     `json:"attributionText,omitempty"`
	BuiltTime          *time.Time   `json:"builtTime,omitempty"`
	ContentIdentifier  []string     `json:"contentIdentifier,omitempty"`
	ContentType        string       `json:"contentType,omitempty"`
	CopyrightText      string       `json:"copyrightText,omitempty"`
	Extension          []string     `json:"extension,omitempty"`
	ExternalIdentifier []string     `json:"externalIdentifier,omitempty"`
	ExternalRef        []string     `json:"externalRef,omitempty"`
	FileKind           string       `json:"fileKind,omitempty"`
	OriginatedBy       []string     `json:"originatedBy,omitempty"`
	ReleaseTime        *time.Time   `json:"releaseTime,omitempty"`
	StandardName       []string     `json:"standardName,omitempty"`
	Summary            string       `json:"summary,omitempty"`
	SuppliedBy         string       `json:"suppliedBy,omitempty"`
	SupportLevel       []string     `json:"supportLevel,omitempty"`
	ValidUntilTime     *time.Time   `json:"validUntilTime,omitempty"`
	VerifiedUsing      []string     `json:"verifiedUsing,omitempty"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, f, &f.PreNode)
}
