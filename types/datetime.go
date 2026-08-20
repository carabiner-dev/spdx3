// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// dateTimeLayout is the lexical form SPDX gives its timestamps. They are
// xsd:dateTimeStamp values, which admit neither fractional seconds nor a
// numeric zone offset, so Go's own time.Time serialization (RFC 3339 with
// nanoseconds and the local offset) is not valid SPDX.
const dateTimeLayout = "2006-01-02T15:04:05Z"

// DateTime is a point in time serialized the way SPDX requires: whole
// seconds, in UTC. Values are converted on the way out, so a time carrying
// a finer precision or another zone still writes a conformant document.
//
// The receivers are deliberately mixed, as they are on Decimal:
// json.Unmarshaler has to take a pointer, while the accessors take values so
// they work on a DateTime that is not addressable.
type DateTime struct {
	Value time.Time
}

// NewDateTime returns a DateTime for the given instant.
func NewDateTime(t time.Time) *DateTime {
	return &DateTime{Value: t}
}

// Time returns the instant this timestamp represents.
func (d DateTime) Time() time.Time {
	return d.Value
}

func (d DateTime) String() string {
	return d.Value.UTC().Format(dateTimeLayout)
}

// UnmarshalJSON reads a timestamp. It accepts any RFC 3339 value, including
// the offsets and sub-second precision SPDX itself does not allow, so that
// documents other tools wrote can still be read.
func (d *DateTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshaling timestamp: %w", err)
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("parsing timestamp %q: %w", s, err)
	}
	d.Value = t
	return nil
}

// MarshalJSON always writes the lexical form SPDX mandates.
func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
