package unmarshal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/carabiner-dev/databom/internal/spdx3/base"
	"github.com/carabiner-dev/databom/internal/spdx3/types"
)

// unmarshal.Node is a universal unmarshaling helper for any type that embeds PreNode.
// It handles both full object serialization and string reference serialization (e.g., "_:id").
func Node(data []byte, target any, preNodePtr *base.PreNode) error {
	// First, check if it's a string reference
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		preNodePtr.ID = strings.TrimPrefix(s, "_:")
		return nil
	}

	// Otherwise, it's an object. Unmarshal into a map first to get all fields
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Use reflection to set fields directly
	return unmarshalFields(raw, target)
}

// unmarshalFields recursively unmarshals fields including embedded structs
func unmarshalFields(raw map[string]json.RawMessage, target any) error {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()

	typeOfNodeSlice := reflect.TypeOf([]types.Node{})
	typeOfNode := reflect.TypeOf(((*types.Node)(nil)))

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if field.Anonymous {
			// Handle embedded structs by recursively unmarshaling
			if fieldValue.CanAddr() {
				if err := unmarshalFields(raw, fieldValue.Addr().Interface()); err != nil {
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
			// If the field is a list of nodes, it can be a list of strings
			// of just the JSONLD ids, or full objects
			if fieldValue.Type() == typeOfNodeSlice {
				if err := unmarshalNodeSlice(rawData, &fieldValue); err != nil {
					return fmt.Errorf("unmarshaling node slice for field %s: %w", tagName, err)
				}
				continue
			} else if fieldValue.Type() == typeOfNode {
				if err := unmarshalNode(rawData, &fieldValue); err != nil {
					return fmt.Errorf("unmarshaling node for field %s: %w", tagName, err)
				}
				continue
			}

			// fmt.Printf("Type: %+v\n", fieldValue.Type().)
			// Create a new value of the field's type
			newVal := reflect.New(fieldValue.Type())
			if err := json.Unmarshal(rawData, newVal.Interface()); err != nil {
				return err
			}
			fieldValue.Set(newVal.Elem())
		}
	}

	return nil
}

func unmarshalNode(rawData json.RawMessage, fieldValue *reflect.Value) error {
	var s string
	if err := json.Unmarshal(rawData, &s); err == nil {
		// It's a string reference, create a NodeRef
		fieldValue.Set(reflect.ValueOf(types.NodeRef{ID: strings.TrimPrefix(s, "_:")}))
		return nil
	}

	node, err := unmarshalNodeDispatch(rawData)
	if err != nil {
		return fmt.Errorf("unmarshaling node: %w", err)
	}
	fieldValue.Set(reflect.ValueOf(node))
	return nil
}

// unmarshalNodeSlice handles unmarshaling a JSON array into a []types.Node slice.
// Each item in the array can be either:
// - A string representing a node reference (e.g., "_:id")
// - A full object that needs to be dispatched to the appropriate concrete type
func unmarshalNodeSlice(rawData json.RawMessage, fieldValue *reflect.Value) error {
	// First, try to unmarshal as a list of raw messages to preserve structure
	var list []json.RawMessage
	if err := json.Unmarshal(rawData, &list); err != nil {
		return fmt.Errorf("unmarshaling as list: %w", err)
	}

	nodeSlice := []types.Node{}
	for i, item := range list {
		// Try to unmarshal as a string first (node reference)
		var s string
		if err := json.Unmarshal(item, &s); err == nil {
			// It's a string reference, create a NodeRef
			nodeSlice = append(nodeSlice, types.NodeRef{ID: strings.TrimPrefix(s, "_:")})
			continue
		}

		// It's not a string, so it must be a full object
		// We need to dispatch it to the appropriate concrete type
		node, err := unmarshalNodeDispatch(item)
		if err != nil {
			return fmt.Errorf("unmarshaling node at index %d: %w", i, err)
		}
		nodeSlice = append(nodeSlice, node)
	}

	fieldValue.Set(reflect.ValueOf(nodeSlice))
	return nil
}

// unmarshalNodeDispatch is a dispatcher function that unmarshals a JSON object
// into the appropriate concrete type. This function is set by the parent package
// to avoid circular dependencies (since profiles import unmarshal, and the
// dispatcher needs to know about profiles).
var unmarshalNodeDispatch func(data []byte) (types.Node, error)

// SetNodeDispatcher sets the function used to dispatch and unmarshal nodes
// based on their type field. This must be called before unmarshaling any
// documents that contain embedded node objects.
func SetNodeDispatcher(dispatcher func(data []byte) (types.Node, error)) {
	unmarshalNodeDispatch = dispatcher
}

func init() {
	// Set a default dispatcher that returns an error if not configured
	unmarshalNodeDispatch = func(data []byte) (types.Node, error) {
		return nil, fmt.Errorf("node dispatching not configured - call unmarshal.SetNodeDispatcher() first")
	}
}
