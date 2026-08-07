// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package marshal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/carabiner-dev/spdx3/types"
)

// NodeMarshaler handles marshaling of SPDX nodes to JSON.
// It implements special logic for nested nodes, serializing them as SPDXID references
// rather than full objects.
type NodeMarshaler struct {
	// ReferenceableIDs lists the identifiers a nested node may be collapsed
	// to a reference string with, normally those of the nodes of the graph
	// being rendered. A nested node identified by anything else has no other
	// place in the document to carry its data, so it is written inline
	// instead of being turned into a reference nothing resolves.
	//
	// A nil map keeps every identifier referenceable, which is the right
	// assumption when marshaling a node on its own, with no graph to
	// consult.
	ReferenceableIDs map[string]struct{}
}

// canReference reports whether a nested node identified by id can be
// collapsed to a reference string.
func (nm *NodeMarshaler) canReference(id string) bool {
	if nm.ReferenceableIDs == nil {
		return true
	}
	_, ok := nm.ReferenceableIDs[id]
	return ok
}

// MarshalNode marshals a node instance to JSON.
// Top-level nodes are serialized fully, but any nested nodes within them
// are serialized as just their SPDXID strings.
func (nm *NodeMarshaler) MarshalNode(source any) ([]byte, error) {
	m, err := nm.marshalToMap(source)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// marshalToMap converts a struct to a map[string]any using reflection,
// serializing the nodes nested in it as references where it can.
func (nm *NodeMarshaler) marshalToMap(source any) (map[string]any, error) {
	result := make(map[string]any)

	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result, nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", v.Kind())
	}

	t := v.Type()
	typeOfNode := reflect.TypeOf((*types.Node)(nil)).Elem()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Handle embedded structs - recursively merge their fields
		if field.Anonymous {
			if fieldValue.Kind() == reflect.Ptr {
				if !fieldValue.IsNil() {
					embedded, err := nm.marshalToMap(fieldValue.Interface())
					if err != nil {
						return nil, fmt.Errorf("marshaling embedded field %s: %w", field.Name, err)
					}
					// Merge embedded fields into result
					for k, v := range embedded {
						result[k] = v
					}
				}
			} else if fieldValue.Kind() == reflect.Struct {
				embedded, err := nm.marshalToMap(fieldValue.Interface())
				if err != nil {
					return nil, fmt.Errorf("marshaling embedded field %s: %w", field.Name, err)
				}
				// Merge embedded fields into result
				for k, v := range embedded {
					result[k] = v
				}
			}
			continue
		}

		// Get JSON tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		parts := strings.Split(jsonTag, ",")
		tagName := parts[0]
		if tagName == "" {
			continue
		}

		hasOmitEmpty := false
		for _, part := range parts[1:] {
			if part == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}

		// Skip zero values for omitempty
		if hasOmitEmpty && isZeroValue(fieldValue) {
			continue
		}

		// Route based on field type
		if fieldValue.Kind() == reflect.Slice {
			elemType := fieldValue.Type().Elem()

			// Check if it's a slice of nodes (interfaces implementing Node)
			if elemType.Kind() == reflect.Interface && elemType.Implements(typeOfNode) {
				marshaled, err := nm.marshalNodeSlice(fieldValue)
				if err != nil {
					return nil, fmt.Errorf("marshaling node slice field %s: %w", tagName, err)
				}
				result[tagName] = marshaled
				continue
			}
		} else if fieldValue.Kind() == reflect.Interface || fieldValue.Kind() == reflect.Ptr {
			// Check if it implements Node
			if fieldValue.Type().Implements(typeOfNode) ||
				(fieldValue.Kind() == reflect.Ptr && fieldValue.Type().Elem().Implements(typeOfNode)) {
				if !fieldValue.IsNil() {
					marshaled, err := nm.marshalSingleNode(fieldValue)
					if err != nil {
						return nil, fmt.Errorf("marshaling node field %s: %w", tagName, err)
					}
					result[tagName] = marshaled
				}
				continue
			}
		}

		// For non-node types, marshal directly
		if !isZeroValue(fieldValue) || !hasOmitEmpty {
			result[tagName] = fieldValue.Interface()
		}
	}

	return result, nil
}

// marshalSingleNode marshals a single node as an SPDXID reference string
func (nm *NodeMarshaler) marshalSingleNode(fieldValue reflect.Value) (any, error) {
	if fieldValue.IsNil() {
		return nil, nil
	}

	node := fieldValue.Interface()

	// Check if it's already a NodeRef - use GetID() since NodeRef.GetSPDXID() returns empty
	if nodeRef, ok := node.(types.NodeRef); ok {
		id := nodeRef.GetID()
		if id == "" {
			return nil, fmt.Errorf("NodeRef has empty ID")
		}
		return id, nil
	}

	// For full nodes, return just the SPDXID as a reference
	if n, ok := node.(types.Node); ok {
		// Try SPDXID first
		id := n.GetSPDXID()
		if id == "" {
			// Fall back to ID if SPDXID is empty (reference-only nodes)
			id = n.GetID()
		}
		if id == "" || !nm.canReference(id) {
			// Nodes without any identifier (hashes, extensions, ...) and
			// nodes whose data is not carried anywhere else in the document
			// can only be represented inline.
			return nm.marshalToMap(node)
		}
		// Return the ID in the same format it was stored
		return id, nil
	}

	return nil, fmt.Errorf("expected Node type, got %T", node)
}

// marshalNodeSlice marshals a slice of nodes as an array of SPDXID reference strings
func (nm *NodeMarshaler) marshalNodeSlice(fieldValue reflect.Value) (any, error) {
	if fieldValue.Len() == 0 {
		return []interface{}{}, nil
	}

	result := make([]interface{}, 0, fieldValue.Len())

	for i := range fieldValue.Len() {
		item := fieldValue.Index(i)

		// Check if it's a NodeRef - use GetID() since NodeRef.GetSPDXID() returns empty
		// This is probably a hack
		if nodeRef, ok := item.Interface().(types.NodeRef); ok {
			id := nodeRef.GetID()
			if id == "" {
				return nil, fmt.Errorf("NodeRef at index %d has empty ID", i)
			}
			result = append(result, id)
			continue
		}

		// For full nodes, extract just the SPDXID
		if node, ok := item.Interface().(types.Node); ok {
			// Try SPDXID first
			id := node.GetSPDXID()
			if id == "" {
				// Fall back to ID if SPDXID is empty (reference-only nodes)
				id = node.GetID()
			}
			if id == "" || !nm.canReference(id) {
				// Nodes without any identifier (hashes, extensions, ...) and
				// nodes whose data is not carried anywhere else in the
				// document can only be represented inline.
				inline, err := nm.marshalToMap(node)
				if err != nil {
					return nil, fmt.Errorf("marshaling inline node at index %d: %w", i, err)
				}
				result = append(result, inline)
				continue
			}
			result = append(result, id)
			continue
		}

		return nil, fmt.Errorf("unexpected type at index %d: %T", i, item.Interface())
	}

	return result, nil
}

// isZeroValue checks if a reflect.Value is the zero value for its type
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Struct:
		// For structs, check if all fields are zero
		for i := range v.NumField() {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}
