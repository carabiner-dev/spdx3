// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package base

// PreNode is a superslim field set that is common to all node types. This is
// effectively embedded in al concrete types and is mainly used for simple
// unmarshaling when detecting node Types
type PreNode struct {
	SPDXID string `json:"spdxID"`
	ID     string `json:"@id"`
	Type   string `json:"type"`
}

func (pn *PreNode) GetSPDXID() string {
	return pn.SPDXID
}

func (pn *PreNode) GetID() string {
	return pn.ID
}

func (pn *PreNode) GetType() string {
	return pn.Type
}
