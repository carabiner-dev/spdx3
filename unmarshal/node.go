// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package unmarshal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/types"
)

// defaultDispatcher backs the Node helper, which the UnmarshalJSON method of
// every node calls. It is only reached when a node is unmarshaled on its own,
// with encoding/json rather than a Parser: the parser drives the traversal
// with its own unmarshaler and never consults this.
var defaultDispatcher types.Dispatcher

// SetDefaultDispatcher sets the package-level dispatcher used by the Node
// helper function. This should be called once during initialization to avoid
// import cycles.
func SetDefaultDispatcher(d types.Dispatcher) {
	defaultDispatcher = d
}

// Node unmarshals JSON data into a target node using the default dispatcher.
// It is what a node's own UnmarshalJSON calls, so that unmarshaling a node
// directly with encoding/json works; the parser does not go through it.
func Node(data []byte, target any, preNodePtr *base.PreNode) error {
	if defaultDispatcher == nil {
		return fmt.Errorf("default dispatcher not set; call unmarshal.SetDefaultDispatcher first")
	}
	nu, ok := defaultDispatcher.(*NodeUnmarshaler)
	if !ok {
		nu = NewNodeUnmarshaler(defaultDispatcher)
	}
	return nu.Unmarshal(data, target, preNodePtr)
}

// Options are the choices a parser makes about how a document is read.
type Options struct {
	// KeepInvalidVocabularyValues stops values that are not members of
	// their vocabulary being dropped as the document is read.
	KeepInvalidVocabularyValues bool
}

// New returns an unmarshaler that resolves node types through classes and
// reads documents according to opts. Everything the unmarshaler needs is
// held on it, so unmarshalers configured differently can run at the same
// time without interfering.
func New(classes map[string]types.Node, opts Options) *NodeUnmarshaler {
	return &NodeUnmarshaler{
		Classes: classes,
		Options: opts,
	}
}

func NewNodeUnmarshaler(dispatcher types.Dispatcher) *NodeUnmarshaler {
	return &NodeUnmarshaler{
		Dispatcher: dispatcher,
	}
}

type NodeUnmarshaler struct {
	// Classes maps an SPDX type name to a prototype of the Go type
	// modelling it. It is the registry an unmarshaler built with New
	// resolves node types through.
	Classes map[string]types.Node

	// Dispatcher resolves node types for an unmarshaler built with
	// NewNodeUnmarshaler, which has no registry of its own.
	Dispatcher types.Dispatcher

	Options Options
}

// preNodeCarrier is satisfied by every node, since they all embed
// base.PreNode, whose GetPreNode method is promoted to them.
type preNodeCarrier interface {
	GetPreNode() *base.PreNode
}

// UnmarshalNode resolves the SPDX type named in the data and reads it into a
// fresh instance of the Go type modelling it. It makes NodeUnmarshaler a
// types.Dispatcher, and recurses through this same unmarshaler rather than
// handing the node back to encoding/json, which would lose the registry and
// the options along the way.
func (nu *NodeUnmarshaler) UnmarshalNode(data []byte) (types.Node, error) {
	if nu.Classes == nil {
		if nu.Dispatcher == nil {
			return nil, fmt.Errorf("unmarshaler has neither a class registry nor a dispatcher")
		}
		return nu.Dispatcher.UnmarshalNode(data)
	}

	prenode := &base.PreNode{}
	if err := json.Unmarshal(data, prenode); err != nil {
		return nil, fmt.Errorf("parsing node: %w", err)
	}

	proto, ok := nu.Classes[prenode.Type]
	if !ok {
		return nil, fmt.Errorf("parsing type %q: %w", prenode.Type, types.ErrUnsupportedNodeType)
	}

	protoType := reflect.TypeOf(proto)
	if protoType.Kind() == reflect.Pointer {
		protoType = protoType.Elem()
	}
	fresh := reflect.New(protoType).Interface()

	node, ok := fresh.(types.Node)
	if !ok {
		return nil, fmt.Errorf("parsing type %q: registered class does not implement types.Node", prenode.Type)
	}
	carrier, ok := fresh.(preNodeCarrier)
	if !ok {
		return nil, fmt.Errorf("parsing type %q: registered class does not embed base.PreNode", prenode.Type)
	}

	if err := nu.Unmarshal(data, node, carrier.GetPreNode()); err != nil {
		return nil, fmt.Errorf("unmarshaling node of type %q: %w", prenode.Type, err)
	}
	return node, nil
}

// unmarshal.NodeUnmarshaler is a universal unmarshaling helper for any type that embeds PreNode.
// It handles both full object serialization and string reference serialization (e.g., "_:id").
func (nu *NodeUnmarshaler) Unmarshal(data []byte, target any, preNodePtr *base.PreNode) error {
	// First, check if it's a string reference. Blank-node labels ("_:id")
	// are kept verbatim so references still match their nodes on output.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		preNodePtr.ID = s
		return nil
	}

	// Otherwise, it's an object. Unmarshal into a map first to get all fields
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Use reflection to set fields directly
	return nu.unmarshalFields(raw, target)
}

// unmarshalFields recursively unmarshals fields including embedded structs
func (nu *NodeUnmarshaler) unmarshalFields(raw map[string]json.RawMessage, target any) error {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	typeOfNode := reflect.TypeOf((*types.Node)(nil)).Elem()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if field.Anonymous {
			// Handle embedded structs by recursively unmarshaling
			if fieldValue.CanAddr() {
				if err := nu.unmarshalFields(raw, fieldValue.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		// Get the JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		// Handle json tag options like "name,omitempty"
		tagName := strings.Split(jsonTag, ",")[0]
		if tagName == "" {
			continue
		}

		// Check if we have data for this field
		if rawData, ok := raw[tagName]; ok {
			// If the field is a slice of nodes (interfaces)
			if fieldValue.Kind() == reflect.Slice && fieldValue.Type().Elem().Kind() == reflect.Interface && fieldValue.Type().Elem().Implements(typeOfNode) {
				if err := nu.unmarshalNodeSlice(rawData, &fieldValue); err != nil {
					return fmt.Errorf("unmarshaling node slice for field %s: %w", tagName, err)
				}
				continue
			}
			// If the field is a single node (interface)
			if fieldValue.Kind() == reflect.Interface && fieldValue.Type().Implements(typeOfNode) {
				if err := nu.unmarshalNode(rawData, &fieldValue); err != nil {
					return fmt.Errorf("unmarshaling node for field %s: %w", tagName, err)
				}
				continue
			}

			newVal := reflect.New(fieldValue.Type())
			if err := nu.unmarshalValue(rawData, newVal.Interface()); err != nil {
				return fmt.Errorf("unmarshaling field %s: %w", tagName, err)
			}
			fieldValue.Set(newVal.Elem())
			nu.dropInvalidVocabularyValues(fieldValue)
		}
	}

	return nil
}

// unmarshalValue reads a field's value. A field holding a node is read
// through this unmarshaler rather than encoding/json, which would route back
// to the node's own UnmarshalJSON and, with it, to the package defaults.
func (nu *NodeUnmarshaler) unmarshalValue(rawData json.RawMessage, target any) error {
	// target is a pointer to the field; a node-valued field is a pointer to
	// a pointer, which has to be allocated before its PreNode can be reached.
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Pointer {
		if bytes.Equal(bytes.TrimSpace(rawData), []byte("null")) {
			return nil
		}
		if v.Elem().IsNil() {
			v.Elem().Set(reflect.New(v.Elem().Type().Elem()))
		}
		if carrier, ok := v.Elem().Interface().(preNodeCarrier); ok {
			return nu.Unmarshal(rawData, v.Elem().Interface(), carrier.GetPreNode())
		}
	}
	if carrier, ok := target.(preNodeCarrier); ok {
		return nu.Unmarshal(rawData, target, carrier.GetPreNode())
	}
	return json.Unmarshal(rawData, target)
}

func (nu *NodeUnmarshaler) unmarshalNode(rawData json.RawMessage, fieldValue *reflect.Value) error {
	var s string
	if err := json.Unmarshal(rawData, &s); err == nil {
		// It's a string reference, create a NodeRef
		fieldValue.Set(reflect.ValueOf(types.NodeRef{ID: s}))
		return nil
	}

	node, err := nu.UnmarshalNode(rawData)
	if err != nil {
		return fmt.Errorf("unmarshaling node: %w", err)
	}
	nodeValue := reflect.ValueOf(node)
	if !nodeValue.Type().AssignableTo(fieldValue.Type()) {
		return fmt.Errorf(
			"%w: node of type %q cannot be assigned to a %s field",
			types.ErrIncompatibleNodeType, node.GetType(), fieldValue.Type(),
		)
	}
	fieldValue.Set(nodeValue)
	return nil
}

// unmarshalNodeSlice handles unmarshaling a JSON array into a slice of types that implement types.Node.
// Each item in the array can be either:
// - A string representing a node reference (e.g., "_:id")
// - A full object that needs to be dispatched to the appropriate concrete type
func (nu *NodeUnmarshaler) unmarshalNodeSlice(rawData json.RawMessage, fieldValue *reflect.Value) error {
	// First, try to unmarshal as a list of raw messages to preserve structure
	var list []json.RawMessage
	if err := json.Unmarshal(rawData, &list); err != nil {
		return fmt.Errorf("unmarshaling as list: %w", err)
	}

	// Create a new slice of the correct type
	sliceType := fieldValue.Type()
	nodeSlice := reflect.MakeSlice(sliceType, 0, len(list))

	for i, item := range list {
		// Try to unmarshal as a string first (node reference)
		var s string
		if err := json.Unmarshal(item, &s); err == nil {
			// It's a string reference, create a NodeRef
			nodeRef := types.NodeRef{ID: s}
			nodeSlice = reflect.Append(nodeSlice, reflect.ValueOf(nodeRef))
			continue
		}

		// It's not a string, so it must be a full object
		// We need to dispatch it to the appropriate concrete type
		node, err := nu.UnmarshalNode(item)
		if err != nil {
			return fmt.Errorf("unmarshaling node at index %d: %w", i, err)
		}
		nodeValue := reflect.ValueOf(node)
		if !nodeValue.Type().AssignableTo(sliceType.Elem()) {
			return fmt.Errorf(
				"%w: node #%d of type %q cannot be assigned to a %s element",
				types.ErrIncompatibleNodeType, i, node.GetType(), sliceType.Elem(),
			)
		}
		nodeSlice = reflect.Append(nodeSlice, nodeValue)
	}

	fieldValue.Set(nodeSlice)
	return nil
}
