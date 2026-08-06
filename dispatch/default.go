// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package dispatch

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/ai"
	"github.com/carabiner-dev/spdx3/profiles/build"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/dataset"
	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/profiles/extension"
	"github.com/carabiner-dev/spdx3/profiles/security"
	"github.com/carabiner-dev/spdx3/profiles/simplelicensing"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

func New() types.Dispatcher {
	classes := map[string]types.Node{}

	// Add the types for all profiles
	maps.Copy(classes, core.Profile.Classes)
	maps.Copy(classes, ai.Profile.Classes)
	maps.Copy(classes, build.Profile.Classes)
	maps.Copy(classes, dataset.Profile.Classes)
	maps.Copy(classes, expandedlicensing.Profile.Classes)
	maps.Copy(classes, extension.Profile.Classes)
	maps.Copy(classes, security.Profile.Classes)
	maps.Copy(classes, simplelicensing.Profile.Classes)
	maps.Copy(classes, software.Profile.Classes)

	return &Default{
		Classes: classes,
	}
}

// Default is the default type dispatcher. Its job is to
// detect its type and pass the correct node type to the
// json unmarshaller.
type Default struct {
	Classes map[string]types.Node
}

// UnmarshalNode takes raw JSON data and unmarshals it into the appropriate
// concrete type based on the "type" field. This is used to handle polymorphic
// node types in SPDX3 JSON-LD documents.
func (d *Default) UnmarshalNode(prenodeData []byte) (types.Node, error) {
	var prenode = &base.PreNode{}
	// Parse the entry to a prenode to determine its type
	if err := json.Unmarshal(prenodeData, prenode); err != nil {
		return nil, fmt.Errorf("parsing node: %w", err)
	}

	proto, ok := d.Classes[prenode.Type]
	if !ok {
		return nil, fmt.Errorf("parsing type %q: %w", prenode.Type, types.ErrUnsupportedNodeType)
	}

	// The registry values are shared prototypes: unmarshal into a fresh
	// instance so that nodes of the same type never alias each other.
	protoType := reflect.TypeOf(proto)
	if protoType.Kind() == reflect.Pointer {
		protoType = protoType.Elem()
	}
	n, ok := reflect.New(protoType).Interface().(types.Node)
	if !ok {
		return nil, fmt.Errorf("parsing type %q: registered class does not implement types.Node", prenode.Type)
	}

	if err := json.Unmarshal(prenodeData, n); err != nil {
		return nil, fmt.Errorf("unmarshaling node of type %q: %w", prenode.Type, err)
	}
	return n, nil
}
