package unmarshal

import (
	"encoding/json"
	"errors"
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
			// of just the JSONLD ids
			if fieldValue.Type() == typeOfNodeSlice {
				list := []any{}
				if err := json.Unmarshal(rawData, &list); err == nil {
					nodeRefs := []types.AddressableById{}
					for _, val := range list {
						if s, ok := val.(string); ok {
							nodeRefs = append(nodeRefs, types.ID(s))
						} else {
							// This needs to be implemented
							//
							// TODO(puerco): Extract the graph dispatcher switch
							return errors.New("unmarshaling embedded node not supported yet")
						}
					}
					nodeSlice := []types.Node{}
					for _, ref := range nodeRefs {
						// Here we should lookup existing nodes
						nodeSlice = append(nodeSlice, types.NodeRef{ID: ref.GetID()})
					}
					fieldValue.Set(reflect.ValueOf(nodeSlice))
					continue
				}
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
