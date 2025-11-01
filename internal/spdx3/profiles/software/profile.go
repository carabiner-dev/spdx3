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
	BuiltTime    *time.Time `json:"builtTime"`
	ReleaseTime  *time.Time `json:"releaseTime"`
	OriginatedBy []string   `json:"originatedBy"`
}

func (f *File) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, f, &f.PreNode)
}
