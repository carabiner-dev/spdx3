// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// The round-trip contract: parsing a document and rendering it back must
// preserve it. Every node must still be findable by its identifier, no
// value may change, no key may appear or vanish (except where omitempty
// legitimately drops an already-empty value), and every reference must
// still resolve. These are the failure modes the library has actually had,
// so they are asserted directly rather than through a golden file, which
// would only pin today's formatting.

const (
	// corpusDir holds the SPDX project's example documents, vendored by
	// ./hack/update-spdx-examples.sh so the tests stay hermetic.
	corpusDir = "testdata/corpus"

	// corpusEnv points the same checks at a different directory, for trying
	// documents that are not vendored here.
	corpusEnv = "SPDX3_CORPUS"
)

func TestRoundTripFixtures(t *testing.T) {
	files, err := filepath.Glob("testdata/roundtrip/*.json")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no round-trip fixtures found")
	files = append(files, "testdata/spdx.json")

	for _, path := range files {
		t.Run(filepath.Base(path), func(t *testing.T) {
			data, err := os.ReadFile(path)
			require.NoError(t, err)
			checkRoundTrip(t, data)
		})
	}
}

// TestRoundTripExamples runs the same contract over the SPDX project's own
// example documents, which exercise shapes no hand-written fixture would
// think to cover.
func TestRoundTripExamples(t *testing.T) {
	dir := os.Getenv(corpusEnv)
	if dir == "" {
		dir = corpusDir
	}

	var checked int
	require.NoError(t, filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".json") {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil || !isSPDX3Document(data) {
			return readErr
		}
		checked++
		rel, relErr := filepath.Rel(dir, path)
		if relErr != nil {
			rel = path
		}
		t.Run(rel, func(t *testing.T) { checkRoundTrip(t, data) })
		return nil
	}))
	require.NotZero(t, checked,
		"no SPDX 3 documents under %s; run ./hack/update-spdx-examples.sh", dir)
	t.Logf("checked %d documents under %s", checked, dir)
}

// isSPDX3Document reports whether data is an SPDX 3 document rather than
// an SPDX 2 one or one of the context and schema files that sit alongside
// them in the example repositories.
func isSPDX3Document(data []byte) bool {
	var probe struct {
		Context json.RawMessage   `json:"@context"`
		Graph   []json.RawMessage `json:"@graph"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return false
	}
	return len(probe.Graph) > 0 && strings.Contains(string(probe.Context), "spdx.org/rdf/3.")
}

// checkRoundTrip parses data, renders it, reparses the result and asserts
// the rendered document says everything the original did.
func checkRoundTrip(t *testing.T, data []byte) {
	t.Helper()

	env, err := NewParser().Parse(bytes.NewReader(data))
	require.NoError(t, err, "parsing the original document")

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf), "rendering the parsed document")

	_, err = NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err, "reparsing the rendered document")

	orig := decodeDocument(t, data)
	rend := decodeDocument(t, buf.Bytes())

	require.Equal(t, orig["@context"], rend["@context"], "@context changed")

	c := &comparer{
		t:          t,
		origTopIDs: nodeIDs(graphNodes(orig)),
		rendTopIDs: nodeIDs(graphNodes(rend)),
		rendAllIDs: map[string]bool{},
	}
	collectIDs(rend, c.rendAllIDs)

	rendByID := map[string]map[string]any{}
	for _, n := range graphNodes(rend) {
		if id := nodeID(n); id != "" {
			rendByID[id] = n
		}
	}

	for _, on := range graphNodes(orig) {
		id := nodeID(on)
		require.NotEmpty(t, id, "original node has no identifier: %v", on)
		rn, ok := rendByID[id]
		if !ok {
			t.Errorf("node %q (%s) is missing from the rendered document", id, nodeType(on))
			continue
		}
		c.compare(nodeType(on), on, rn)
	}

	// Every reference in the rendered document must resolve to a node the
	// document actually contains.
	scanStrings(rend, func(s string) {
		if !strings.HasPrefix(s, "_:") && !c.origTopIDs[s] {
			return
		}
		if !c.rendAllIDs[s] {
			t.Errorf("reference %q resolves to nothing in the rendered document", s)
		}
	})
}

type comparer struct {
	t                      *testing.T
	origTopIDs, rendTopIDs map[string]bool
	rendAllIDs             map[string]bool
}

// compare walks an original value and its rendered counterpart. Array
// indices are left out of the reported path so failures read as the
// property that broke.
func (c *comparer) compare(path string, orig, rend any) {
	// A node may legitimately be written either inline or as a reference to
	// it, as long as its data is still somewhere in the document.
	if id, ok := orig.(string); ok {
		if obj, isObj := rend.(map[string]any); isObj && nodeID(obj) == id {
			return
		}
	}
	if id, ok := rend.(string); ok {
		if obj, isObj := orig.(map[string]any); isObj && nodeID(obj) == id {
			if !c.rendTopIDs[id] {
				c.t.Errorf("%s: node %q was inlined in the original but rendered as a reference to a node the document does not contain", path, id)
			}
			return
		}
	}

	switch ov := orig.(type) {
	case map[string]any:
		rv, ok := rend.(map[string]any)
		if !ok {
			c.t.Errorf("%s: expected an object, rendered %s", path, describe(rend))
			return
		}
		for k, val := range ov {
			sub := path + "." + k
			rval, present := rv[k]
			if !present {
				// omitempty legitimately drops values that carry nothing.
				if !isEmptyValue(val) {
					c.t.Errorf("%s: dropped, original value was %s", sub, describe(val))
				}
				continue
			}
			c.compare(sub, val, rval)
		}
		for k := range rv {
			if _, present := ov[k]; !present {
				c.t.Errorf("%s.%s: appeared in the rendered document", path, k)
			}
		}
	case []any:
		rv, ok := rend.([]any)
		if !ok {
			c.t.Errorf("%s: expected an array, rendered %s", path, describe(rend))
			return
		}
		if len(ov) != len(rv) {
			c.t.Errorf("%s: had %d entries, rendered %d", path, len(ov), len(rv))
			return
		}
		for i := range ov {
			c.compare(path, ov[i], rv[i])
		}
	default:
		if !leafEqual(orig, rend) {
			c.t.Errorf("%s: was %s, rendered %s", path, describe(orig), describe(rend))
		}
	}
}

// leafEqual compares scalars. JSON-LD carries some numeric datatypes as
// strings to preserve their RDF type, so a number and its exact lexical
// form are the same value; a different lexical form is not.
func leafEqual(a, b any) bool {
	if an, ok := a.(json.Number); ok {
		if bs, ok := b.(string); ok {
			return an.String() == bs
		}
	}
	if as, ok := a.(string); ok {
		if bn, ok := b.(json.Number); ok {
			return as == bn.String()
		}
	}
	return a == b
}

// isEmptyValue reports whether a value carries no information, and may
// therefore be dropped by omitempty. Note that false and 0 do carry
// information: dropping them loses an explicit assertion.
func isEmptyValue(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return t == ""
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	}
	return false
}

func describe(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	if len(data) > 160 {
		return string(data[:157]) + "..."
	}
	return string(data)
}

func decodeDocument(t *testing.T, data []byte) map[string]any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	doc := map[string]any{}
	require.NoError(t, dec.Decode(&doc))
	return doc
}

func graphNodes(doc map[string]any) []map[string]any {
	raw, ok := doc["@graph"].([]any)
	if !ok {
		return nil
	}
	nodes := make([]map[string]any, 0, len(raw))
	for _, n := range raw {
		if m, ok := n.(map[string]any); ok {
			nodes = append(nodes, m)
		}
	}
	return nodes
}

// nodeID returns the identifier a node is referenced by: its spdxId, or
// its @id for the blank nodes that carry no spdxId.
func nodeID(m map[string]any) string {
	if s, ok := m["spdxId"].(string); ok && s != "" {
		return s
	}
	if s, ok := m["@id"].(string); ok {
		return s
	}
	return ""
}

func nodeType(m map[string]any) string {
	if s, ok := m["type"].(string); ok {
		return s
	}
	return "?"
}

func nodeIDs(nodes []map[string]any) map[string]bool {
	ids := map[string]bool{}
	for _, n := range nodes {
		if id := nodeID(n); id != "" {
			ids[id] = true
		}
	}
	return ids
}

// collectIDs gathers every identifier appearing anywhere in a document,
// including those of inlined nodes.
func collectIDs(v any, ids map[string]bool) {
	switch t := v.(type) {
	case map[string]any:
		for _, key := range []string{"spdxId", "@id"} {
			if s, ok := t[key].(string); ok && s != "" {
				ids[s] = true
			}
		}
		for _, sub := range t {
			collectIDs(sub, ids)
		}
	case []any:
		for _, sub := range t {
			collectIDs(sub, ids)
		}
	}
}

func scanStrings(v any, fn func(string)) {
	switch t := v.(type) {
	case string:
		fn(t)
	case map[string]any:
		for _, sub := range t {
			scanStrings(sub, fn)
		}
	case []any:
		for _, sub := range t {
			scanStrings(sub, fn)
		}
	}
}
