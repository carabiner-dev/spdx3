// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// canonicalTime is the lexical form SPDX asks for, as JSON.
const canonicalTime = `"2024-05-31T00:00:00Z"`

func TestDateTime(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    string
		expected string
	}{
		{"canonical form is kept", canonicalTime, canonicalTime},
		{"fractional seconds are dropped", `"2024-05-31T00:00:00.123456789Z"`, canonicalTime},
		{"offsets are converted to UTC", `"2024-05-31T00:00:00-06:00"`, `"2024-05-31T06:00:00Z"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var d DateTime
			require.NoError(t, json.Unmarshal([]byte(tc.input), &d))
			out, err := json.Marshal(d)
			require.NoError(t, err)
			require.JSONEq(t, tc.expected, string(out))
		})
	}

	t.Run("rejects values that are not timestamps", func(t *testing.T) {
		var d DateTime
		require.Error(t, json.Unmarshal([]byte(`"not-a-time"`), &d))
		require.Error(t, json.Unmarshal([]byte(`17`), &d))
	})

	// The failure this type exists to prevent: Go's own time serialization
	// carries nanoseconds and the local offset, neither of which SPDX allows.
	t.Run("a local time with nanoseconds still writes a valid value", func(t *testing.T) {
		zone := time.FixedZone("UTC-6", -6*60*60)
		out, err := json.Marshal(NewDateTime(time.Date(2026, 8, 6, 18, 30, 15, 123456789, zone)))
		require.NoError(t, err)
		require.JSONEq(t, `"2026-08-07T00:30:15Z"`, string(out))
	})
}
