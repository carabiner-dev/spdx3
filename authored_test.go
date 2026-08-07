// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"bytes"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
)

// dateTimeStamp is the lexical form the SPDX 3.0.1 schema requires of
// timestamps: whole seconds, UTC, no offset.
var dateTimeStamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)

// looksLikeTimestamp catches anything date-shaped, so a value in the wrong
// lexical form is reported rather than skipped.
var looksLikeTimestamp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)

// TestAuthoredDocumentTimestamps renders a document built in Go rather than
// read from a file. The round-trip tests cannot catch a bad timestamp, since
// every timestamp in the corpus is already in the form the spec wants and is
// therefore written back unchanged; only a document we author ourselves puts
// a Go time.Time through the marshaler.
func TestAuthoredDocumentTimestamps(t *testing.T) {
	// A real time.Now(): sub-second precision, and whatever zone the machine
	// running the tests happens to be in.
	now := time.Now()
	require.NotZero(t, now.Nanosecond(), "this test needs a time carrying nanoseconds")

	creation := &core.CreationInfo{
		PreNode:     base.PreNode{ID: creationInfoID, Type: core.CreationInfoClass},
		SpecVersion: specVersion301,
		Created:     types.NewDateTime(now),
		CreatedBy:   []core.AgentDescendant{types.NodeRef{ID: "https://example.com/spdx/alice"}},
	}

	pkg := &software.Package{
		SoftwareArtifact: software.SoftwareArtifact{
			Artifact: core.Artifact{
				Node: core.Node{
					PreNode:      base.PreNode{SPDXID: "https://example.com/spdx/pkg", Type: "software_Package"},
					Name:         "example-lib",
					CreationInfo: creation,
				},
				BuiltTime:   types.NewDateTime(now.Add(-time.Hour)),
				ReleaseTime: types.NewDateTime(now.UTC()),
			},
		},
		PackageVersion: packageVersion,
	}

	env := &Envelope{
		Context: NewContext(ContextURL301),
		Graph:   Graph{creation, pkg},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, (&Renderer{}).Render(env, buf))

	doc := decodeDocument(t, buf.Bytes())
	var checked int
	scanStrings(doc, func(s string) {
		if !looksLikeTimestamp.MatchString(s) {
			return
		}
		checked++
		require.Regexp(t, dateTimeStamp, s,
			"a timestamp the library wrote is not in the form SPDX requires")
	})
	require.Equal(t, 3, checked, "expected to have checked created, builtTime and releaseTime")

	// And the document still reads back.
	reparsed, err := NewParser().Parse(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	require.Len(t, reparsed.Graph, 2)
}
