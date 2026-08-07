// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/carabiner-dev/spdx3/types"
)

// Finding is a place where a document departs from the specification.
type Finding struct {
	// NodeID identifies the node the finding is in, by its spdxId or, for
	// blank nodes, its @id.
	NodeID string
	// NodeType is the SPDX class of that node.
	NodeType string
	// Property is the JSON name of the property at fault.
	Property string
	// Value is the offending value.
	Value string
	// Message says what is wrong with it.
	Message string
}

func (f Finding) String() string {
	where := f.NodeID
	if where == "" {
		where = "(unidentified node)"
	}
	return fmt.Sprintf("%s (%s): %s: %s", where, f.NodeType, f.Property, f.Message)
}

// Validate reports the values in a document that are not members of the
// vocabulary their property is drawn from.
//
// A parser drops such values as it reads, so validating a document parsed
// the usual way finds nothing. Parse it with WithInvalidVocabularyValues to
// keep them, and validate that.
//
// This checks vocabularies only. It says nothing about required properties,
// cardinalities or the ranges of object properties.
func Validate(env *Envelope) []Finding {
	if env == nil {
		return nil
	}
	findings := []Finding{}
	for _, node := range env.Graph {
		if node == nil {
			continue
		}
		id := node.GetSPDXID()
		if id == "" {
			id = node.GetID()
		}
		findings = append(findings, validateNode(reflect.ValueOf(node), id, node.GetType())...)
	}
	return findings
}

var vocabularyValueType = reflect.TypeOf((*types.VocabularyValue)(nil)).Elem()

// validateNode walks the exported fields of a node, descending through the
// structs embedded in it, and reports every vocabulary value it does not
// recognize.
func validateNode(v reflect.Value, id, nodeType string) []Finding {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	findings := []Finding{}
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		value := v.Field(i)

		if field.Anonymous {
			findings = append(findings, validateNode(value, id, nodeType)...)
			continue
		}

		property := strings.Split(field.Tag.Get("json"), ",")[0]
		if property == "" || property == "-" {
			continue
		}

		switch {
		case value.Type().Implements(vocabularyValueType):
			if f, bad := checkVocabularyValue(value, id, nodeType, property); bad {
				findings = append(findings, f)
			}
		case value.Kind() == reflect.Slice && value.Type().Elem().Implements(vocabularyValueType):
			for j := range value.Len() {
				if f, bad := checkVocabularyValue(value.Index(j), id, nodeType, property); bad {
					findings = append(findings, f)
				}
			}
		}
	}
	return findings
}

func checkVocabularyValue(v reflect.Value, id, nodeType, property string) (Finding, bool) {
	if v.IsZero() {
		return Finding{}, false
	}
	value, ok := v.Interface().(types.VocabularyValue)
	if !ok || value.IsValid() {
		return Finding{}, false
	}
	return Finding{
		NodeID:   id,
		NodeType: nodeType,
		Property: property,
		Value:    fmt.Sprintf("%v", v.Interface()),
		Message:  fmt.Sprintf("%q is not a member of the vocabulary this property draws from", v.Interface()),
	}, true
}
