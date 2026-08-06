// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecimal(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected float64
		output   string
	}{
		{"string form keeps lexical", `"0.0000036"`, 0.0000036, `"0.0000036"`},
		{"string form keeps trailing zero", `"5.0"`, 5.0, `"5.0"`},
		{"number form", `4.3`, 4.3, `"4.3"`},
		{"integer number form", `7`, 7, `"7"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var d Decimal
			require.NoError(t, json.Unmarshal([]byte(tc.input), &d))
			require.Equal(t, tc.expected, d.Value)
			out, err := json.Marshal(d)
			require.NoError(t, err)
			require.Equal(t, tc.output, string(out))
		})
	}

	var d Decimal
	require.Error(t, json.Unmarshal([]byte(`"not-a-number"`), &d))
}
