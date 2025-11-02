package build

import (
	"time"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const Prefix = "build"

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		"Build": &Build{},
	},
}

// Build represents the act of converting software inputs into software artifacts
type Build struct {
	core.Element
	BuildEndTime           *time.Time             `json:"buildEndTime,omitempty"`
	BuildId                string                 `json:"buildId,omitempty"`
	BuildStartTime         *time.Time             `json:"buildStartTime,omitempty"`
	BuildType              string                 `json:"buildType"`
	ConfigSourceDigest     []core.Hash            `json:"configSourceDigest,omitempty"`
	ConfigSourceEntrypoint []string               `json:"configSourceEntrypoint,omitempty"`
	ConfigSourceUri        []string               `json:"configSourceUri,omitempty"`
	Environment            []core.DictionaryEntry `json:"environment,omitempty"`
	Parameter              []core.DictionaryEntry `json:"parameter,omitempty"`
}

func (b *Build) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}
