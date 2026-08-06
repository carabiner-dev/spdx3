// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

type CreationInfo struct {
	base.PreNode
	root         bool              `json:""`
	Name         string            `json:"name,omitempty"`
	SpecVersion  string            `json:"specVersion"`
	CreatedBy    []AgentDescendant `json:"createdBy"`
	CreatedUsing []types.Node      `json:"createdUsing,omitempty"`
	Created      *time.Time        `json:"created"`
	Comment      string            `json:"comment,omitempty"`
}

func (ci *CreationInfo) GetType() string {
	return CreationInfoClass
}

func (ci *CreationInfo) GetName() string {
	return ci.Name
}

func (ci *CreationInfo) GetCreationInfo() types.Node {
	return ci
}

// func (ci *CreationInfo) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(ci)

// 	// if ci.root {
// 	// 	return json.Marshal(ci)
// 	// }

// 	// return fmt.Appendf(nil, "\"_:%s\"", ci.ID), nil
// }

func (ci *CreationInfo) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ci, &ci.PreNode)
}
