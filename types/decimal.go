// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Decimal represents an xsd:decimal value. The SPDX JSON-LD serialization
// carries decimals as strings to preserve the RDF datatype (JSON numbers
// would become xsd:integer/xsd:double), but this type also accepts plain
// JSON numbers on input.
type Decimal struct {
	Value float64
	raw   string
}

func NewDecimal(v float64) Decimal {
	return Decimal{Value: v}
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return fmt.Errorf("parsing decimal %q: %w", s, err)
		}
		d.Value = v
		d.raw = s
		return nil
	}
	d.raw = ""
	return json.Unmarshal(data, &d.Value)
}

// MarshalJSON always emits the string lexical form, preserving the exact
// input representation when there was one.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.raw != "" {
		return json.Marshal(d.raw)
	}
	return json.Marshal(strconv.FormatFloat(d.Value, 'f', -1, 64))
}
