// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"reflect"

	"github.com/carabiner-dev/spdx3/types"
)

var typeOfNode = reflect.TypeOf((*types.Node)(nil)).Elem()

// resolveReferences replaces the reference stubs left by parsing with the
// nodes they name, so that following a property leads to the element itself
// rather than to its identifier. A reference the document does not carry a
// node for — one naming an element in another document, or an individual the
// specification predefines — is left as it was.
//
// Rendering is unaffected: a nested node is written as a reference when the
// graph carries it, which is precisely the case where one was resolved.
func (g Graph) resolveReferences() {
	index := make(map[string]types.Node, len(g)*2)
	for _, node := range g {
		if node == nil {
			continue
		}
		for _, id := range []string{node.GetSPDXID(), node.GetID()} {
			if id != "" {
				index[id] = node
			}
		}
	}
	if len(index) == 0 {
		return
	}

	visited := map[uintptr]bool{}
	for _, node := range g {
		resolveNode(reflect.ValueOf(node), index, visited)
	}
}

// resolveNode walks one node, once. A node already visited is skipped, so a
// graph whose references form a cycle terminates.
func resolveNode(v reflect.Value, index map[string]types.Node, visited map[uintptr]bool) {
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return
	}
	if visited[v.Pointer()] {
		return
	}
	visited[v.Pointer()] = true
	walkFields(v.Elem(), index, visited)
}

// walkFields resolves the node-valued properties of a struct. Embedded
// structs are walked as part of the same node rather than as nodes of their
// own: one embedded at the start of another shares its address, so treating
// it as a node would make it look already visited.
func walkFields(v reflect.Value, index map[string]types.Node, visited map[uintptr]bool) {
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		value := v.Field(i)

		if field.Anonymous {
			embedded := value
			if embedded.Kind() == reflect.Pointer {
				if embedded.IsNil() {
					continue
				}
				embedded = embedded.Elem()
			}
			walkFields(embedded, index, visited)
			continue
		}

		switch {
		case holdsNode(value.Type()):
			resolveValue(value, index, visited)
		case value.Kind() == reflect.Slice && holdsNode(value.Type().Elem()):
			for j := range value.Len() {
				resolveValue(value.Index(j), index, visited)
			}
		}
	}
}

// holdsNode reports whether a field of this type carries a node: either an
// interface a node satisfies, or a pointer to a node's struct.
func holdsNode(t reflect.Type) bool {
	if t.Kind() == reflect.Interface {
		return t.Implements(typeOfNode) || typeOfNode.Implements(t)
	}
	return t.Kind() == reflect.Pointer && t.Implements(typeOfNode)
}

// resolveValue points one field at the node it names, when the graph has it.
func resolveValue(value reflect.Value, index map[string]types.Node, visited map[uintptr]bool) {
	if !value.IsValid() || (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) && value.IsNil() {
		return
	}
	current, ok := value.Interface().(types.Node)
	if !ok {
		return
	}

	if target, found := index[referencedID(current)]; found && value.CanSet() {
		targetValue := reflect.ValueOf(target)
		// A reference whose target is not valid here is left alone: the
		// document is wrong, but that is Validate's business, not this pass's.
		if targetValue.Type().AssignableTo(value.Type()) && targetValue.Interface() != current {
			value.Set(targetValue)
			// The target is a node of the graph and is walked in its own right.
			return
		}
	}

	// Not a reference to anything the graph holds: it may be a node written
	// inline, whose own references still need resolving.
	resolveNode(reflect.ValueOf(current), index, visited)
}

// referencedID returns the identifier a value names. A node states its own
// with spdxId, while a reference to one carries it as the node's @id.
func referencedID(n types.Node) string {
	if id := n.GetSPDXID(); id != "" {
		return id
	}
	return n.GetID()
}
