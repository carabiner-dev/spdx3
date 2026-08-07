// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package unmarshal

import (
	"reflect"

	"github.com/carabiner-dev/spdx3/types"
)

// keepInvalidVocabularyValues turns off the check that drops values outside
// their vocabulary while parsing. It is package level for the same reason
// the default dispatcher is: the parser has no other way to reach the
// reflective unmarshaler the profiles call into.
var keepInvalidVocabularyValues bool

// SetKeepInvalidVocabularyValues controls what happens to a property whose
// value is not a member of its vocabulary. By default such a value is
// dropped while parsing, leaving the property unset. Passing true keeps it,
// so the document is read exactly as written and Validate can report on it.
func SetKeepInvalidVocabularyValues(keep bool) {
	keepInvalidVocabularyValues = keep
}

var vocabularyValueType = reflect.TypeOf((*types.VocabularyValue)(nil)).Elem()

// dropInvalidVocabularyValues clears a field holding a value outside its
// vocabulary, and removes such entries from a field holding a list of them.
// Anything that is not a vocabulary value is left alone.
func dropInvalidVocabularyValues(field reflect.Value) {
	if keepInvalidVocabularyValues || !field.IsValid() {
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
