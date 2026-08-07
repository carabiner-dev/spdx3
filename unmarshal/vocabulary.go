// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package unmarshal

import (
	"reflect"

	"github.com/carabiner-dev/spdx3/types"
)

var vocabularyValueType = reflect.TypeOf((*types.VocabularyValue)(nil)).Elem()

// dropInvalidVocabularyValues clears a field holding a value outside its
// vocabulary, and removes such entries from a field holding a list of them.
// Anything that is not a vocabulary value is left alone.
func (nu *NodeUnmarshaler) dropInvalidVocabularyValues(field reflect.Value) {
	if nu.Options.KeepInvalidVocabularyValues || !field.IsValid() {
		return
	}

	if field.Type().Implements(vocabularyValueType) {
		if !isValidVocabularyValue(field) {
			field.SetZero()
		}
		return
	}

	if field.Kind() != reflect.Slice || !field.Type().Elem().Implements(vocabularyValueType) {
		return
	}
	kept := reflect.MakeSlice(field.Type(), 0, field.Len())
	for i := range field.Len() {
		if isValidVocabularyValue(field.Index(i)) {
			kept = reflect.Append(kept, field.Index(i))
		}
	}
	field.Set(kept)
}

// isValidVocabularyValue reports whether a value belongs to its vocabulary.
// An unset value is left for the absent-property rules to deal with rather
// than being treated as a bad one.
func isValidVocabularyValue(v reflect.Value) bool {
	if v.IsZero() {
		return true
	}
	value, ok := v.Interface().(types.VocabularyValue)
	return !ok || value.IsValid()
}
