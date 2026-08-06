// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

// CreationInfo records who created an Element and when. It is not an
// Element itself, so it carries none of Element's properties: the model
// gives it only the five below.
type CreationInfo struct {
	base.PreNode
	SpecVersion  string            `json:"specVersion"`
	CreatedBy    []AgentDescendant `json:"createdBy"`
	CreatedUsing []types.Node      `json:"createdUsing,omitempty"`
	Created      *time.Time        `json:"created"`
	Comment      string            `json:"comment,omitempty"`
}

func (ci *CreationInfo) GetType() string {
	return CreationInfoClass
}

// GetName satisfies types.Node. CreationInfo has no name in the model.
func (ci *CreationInfo) GetName() string {
	return ""
}

func (ci *CreationInfo) GetCreationInfo() types.Node {
	return ci
}

func (ci *CreationInfo) UnmarshalJSON(data []byte) error {
	return unmarshal.Node(data, ci, &ci.PreNode)
}
