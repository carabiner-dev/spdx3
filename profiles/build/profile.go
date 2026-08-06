// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"time"

	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

const (
	Prefix    = "build"
	BuildType = "Build"
)

var Profile = types.Profile{
	Prefix: Prefix,
	Classes: map[string]types.Node{
		BuildType:                               &Build{},
		fmt.Sprintf("%s_%s", Prefix, BuildType): &Build{},
	},
}

// Build represents the act of converting software inputs into software artifacts
type Build struct {
	core.Element
	BuildEndTime           *time.Time             `json:"build_buildEndTime,omitempty"`
	BuildId                string                 `json:"build_buildId,omitempty"`
	BuildStartTime         *time.Time             `json:"build_buildStartTime,omitempty"`
	BuildType              string                 `json:"build_buildType"`
	ConfigSourceDigest     []core.Hash            `json:"build_configSourceDigest,omitempty"`
	ConfigSourceEntrypoint []string               `json:"build_configSourceEntrypoint,omitempty"`
	ConfigSourceUri        []string               `json:"build_configSourceUri,omitempty"`
	Environment            []core.DictionaryEntry `json:"build_environment,omitempty"`
	Parameter              []core.DictionaryEntry `json:"build_parameter,omitempty"`
}

func (b *Build) GetType() string {
	return fmt.Sprintf("%s_%s", Prefix, BuildType)
}

func (b *Build) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, b, &b.PreNode)
}
