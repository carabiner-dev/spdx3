// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

// ContextURL301 is the JSON-LD context URL of SPDX 3.0.1, the spec version
// this library targets.
const ContextURL301 = "https://spdx.org/rdf/3.0.1/spdx-context.jsonld"

// contextVersionPatterns match the context URLs the SPDX project publishes:
// spdx.org/rdf/<version>/spdx-context.jsonld for releases and
// spdx.github.io/spdx-spec/v<version>/rdf/spdx-context.jsonld for release
// candidates and development builds.
var contextVersionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^https?://spdx\.org/rdf/([0-9][^/]*)/spdx-context\.jsonld$`),
	regexp.MustCompile(`^https?://spdx\.github\.io/spdx-spec/v([^/]+)/rdf/spdx-context\.jsonld$`),
}

// Context holds a document's JSON-LD @context. SPDX 3 documents normally
// carry a single context URL string, but JSON-LD also admits array and
// object forms; those parse too and are preserved verbatim when rendering.
// The receivers are deliberately mixed: json.Unmarshaler has to take a
// pointer, while the accessors take values so they work on the Context field
// of an Envelope value, which is not addressable when the envelope is
// rendered.
type Context struct {
	url string
	raw json.RawMessage
}

// NewContext returns a Context referencing a single context URL, the form
// the SPDX serialization mandates.
func NewContext(url string) Context {
	return Context{url: url}
}

func (c *Context) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	var s string
	if err := json.Unmarshal(data, &s); err == nil && !bytes.Equal(data, []byte("null")) {
		c.url, c.raw = s, nil
		return nil
	}
	if len(data) == 0 || (data[0] != '[' && data[0] != '{') {
		return fmt.Errorf("@context must be a string, array or object")
	}
	c.url = ""
	c.raw = append(json.RawMessage(nil), data...)
	return nil
}

func (c Context) MarshalJSON() ([]byte, error) {
	if c.raw != nil {
		return c.raw, nil
	}
	return json.Marshal(c.url)
}

// String returns the context URL or, for the array and object forms, the
// raw JSON text.
func (c Context) String() string {
	if c.raw != nil {
		return string(c.raw)
	}
	return c.url
}

// URLs returns the context URLs the @context value references: the URL
// itself in the common string form or every string entry of the array form.
// Inline object contexts define their terms directly and reference none.
func (c Context) URLs() []string {
	if c.url != "" {
		return []string{c.url}
	}
	if len(c.raw) == 0 || c.raw[0] != '[' {
		return nil
	}
	items := []json.RawMessage{}
	if err := json.Unmarshal(c.raw, &items); err != nil {
		return nil
	}
	urls := []string{}
	for _, item := range items {
		var s string
		if err := json.Unmarshal(item, &s); err == nil {
			urls = append(urls, s)
		}
	}
	return urls
}

// Version returns the SPDX spec version the context URL pins ("3.0.1",
// "3.1", "3.1-RC1", ...) or an empty string when the context references no
// recognizable SPDX context URL.
func (c Context) Version() string {
	for _, u := range c.URLs() {
		for _, re := range contextVersionPatterns {
			if m := re.FindStringSubmatch(u); m != nil {
				return m[1]
			}
		}
	}
	return ""
}
