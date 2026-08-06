// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package marshal

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/stretchr/testify/require"
)

func TestNodeMarshaler_MarshalNode(t *testing.T) {
	marshaler := &NodeMarshaler{}

	t.Run("Simple Person node", func(t *testing.T) {
		person := &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						ID:     "_:person1",
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name: "John Doe",
				},
			},
		}

		data, err := marshaler.MarshalNode(person)
		require.NoError(t, err)

		// Unmarshal back to verify structure
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "_:person1", result["@id"])
		require.Equal(t, "Person", result["type"])
		require.Equal(t, "SPDXRef-Person1", result["spdxId"])
		require.Equal(t, "John Doe", result["name"])
	})

	t.Run("CreationInfo with nested node reference", func(t *testing.T) {
		now := time.Now()
		person := &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						ID:     "https://spdx.org/spdxdocs/Person1-1000e6a2-0229-4875-baa7-c99be213b6e1",
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name: "John Doe",
				},
			},
		}

		creationInfo := &core.CreationInfo{
			PreNode: base.PreNode{
				ID:     "_:creationinfo",
				Type:   "CreationInfo",
				SPDXID: "SPDXRef-CreationInfo",
			},
			Name:        "Test Creation",
			SpecVersion: "3.0.1",
			CreatedBy:   []core.AgentDescendant{person},
			Created:     &now,
		}

		data, err := marshaler.MarshalNode(creationInfo)
		require.NoError(t, err)

		// Unmarshal back to verify structure
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "_:creationinfo", result["@id"])
		require.Equal(t, "CreationInfo", result["type"])
		require.Equal(t, "3.0.1", result["specVersion"])

		// Check that CreatedBy is serialized as an array of SPDXID strings
		createdBy, ok := result["createdBy"].([]interface{})
		require.True(t, ok, "createdBy should be an array")
		require.Len(t, createdBy, 1)
		require.Equal(t, "SPDXRef-Person1", createdBy[0])
	})

	t.Run("Node with NodeRef in CreatedBy", func(t *testing.T) {
		now := time.Now()
		creationInfo := &core.CreationInfo{
			PreNode: base.PreNode{
				ID:     "_:creationinfo",
				Type:   "CreationInfo",
				SPDXID: "SPDXRef-CreationInfo",
			},
			Name:        "Test Creation",
			SpecVersion: "3.0.1",
			CreatedBy: []core.AgentDescendant{
				types.NodeRef{ID: "https://spdx.org/spdxdocs/Person1-abc123"},
			},
			Created: &now,
		}

		data, err := marshaler.MarshalNode(creationInfo)
		require.NoError(t, err)

		// Unmarshal back to verify structure
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Check that CreatedBy is serialized as an array of strings
		// NodeRef uses GetID() which returns the ID field value
		createdBy, ok := result["createdBy"].([]interface{})
		require.True(t, ok, "createdBy should be an array")
		require.Len(t, createdBy, 1)
		require.Equal(t, "https://spdx.org/spdxdocs/Person1-abc123", createdBy[0])
	})

	t.Run("Person with nested CreationInfo reference", func(t *testing.T) {
		creationInfo := &core.CreationInfo{
			PreNode: base.PreNode{
				ID:     "_:creationinfo",
				Type:   "CreationInfo",
				SPDXID: "SPDXRef-CreationInfo",
			},
		}

		person := &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						ID:     "_:person1",
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name:         "John Doe",
					CreationInfo: creationInfo,
				},
			},
		}

		data, err := marshaler.MarshalNode(person)
		require.NoError(t, err)

		// Unmarshal back to verify structure
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		require.Equal(t, "_:person1", result["@id"])
		require.Equal(t, "Person", result["type"])
		require.Equal(t, "John Doe", result["name"])

		// Check that creationInfo is serialized as a SPDXID reference string
		creationInfoID, ok := result["creationInfo"].(string)
		require.True(t, ok, "creationInfo should be a string")
		require.Equal(t, "SPDXRef-CreationInfo", creationInfoID)
	})

	t.Run("Omitempty fields", func(t *testing.T) {
		person := &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						ID:     "_:person1",
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name: "John Doe",
					// Comment is omitted (zero value)
				},
			},
		}

		data, err := marshaler.MarshalNode(person)
		require.NoError(t, err)

		// Unmarshal back to verify structure
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Comment should not be present
		_, hasComment := result["comment"]
		require.False(t, hasComment, "comment should be omitted when empty")
	})
}

func TestReferenceableIDs(t *testing.T) {
	newPerson := func() *core.Person {
		return &core.Person{
			Agent: core.Agent{
				Node: core.Node{
					PreNode: base.PreNode{
						Type:   "Person",
						SPDXID: "SPDXRef-Person1",
					},
					Name: "John Doe",
				},
			},
		}
	}

	marshalCreatedBy := func(t *testing.T, nm *NodeMarshaler) any {
		t.Helper()
		data, err := nm.MarshalNode(&core.CreationInfo{
			PreNode:   base.PreNode{ID: "_:creationinfo", Type: "CreationInfo"},
			CreatedBy: []core.AgentDescendant{newPerson()},
		})
		require.NoError(t, err)

		var result map[string]any
		require.NoError(t, json.Unmarshal(data, &result))
		createdBy, ok := result["createdBy"].([]any)
		require.True(t, ok, "createdBy should be an array")
		require.Len(t, createdBy, 1)
		return createdBy[0]
	}

	t.Run("nil map keeps every id referenceable", func(t *testing.T) {
		require.Equal(t, "SPDXRef-Person1", marshalCreatedBy(t, &NodeMarshaler{}))
	})

	t.Run("known id collapses to a reference", func(t *testing.T) {
		nm := &NodeMarshaler{ReferenceableIDs: map[string]struct{}{
			"SPDXRef-Person1": {},
		}}
		require.Equal(t, "SPDXRef-Person1", marshalCreatedBy(t, nm))
	})

	t.Run("unknown id is written inline", func(t *testing.T) {
		nm := &NodeMarshaler{ReferenceableIDs: map[string]struct{}{
			"SPDXRef-SomethingElse": {},
		}}
		inline, ok := marshalCreatedBy(t, nm).(map[string]any)
		require.True(t, ok, "an unreferenceable node should be inlined")
		require.Equal(t, "SPDXRef-Person1", inline["spdxId"])
		require.Equal(t, "John Doe", inline["name"])
	})
}

func TestIsZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		getValue func() interface{}
		want     bool
	}{
		{
			name:     "empty string",
			getValue: func() interface{} { return "" },
			want:     true,
		},
		{
			name:     "non-empty string",
			getValue: func() interface{} { return "test" },
			want:     false,
		},
		{
			name:     "nil pointer",
			getValue: func() interface{} { var p *int; return p },
			want:     true,
		},
		{
			name:     "non-nil pointer",
			getValue: func() interface{} { v := 42; return &v },
			want:     false,
		},
		{
			name:     "empty slice",
			getValue: func() interface{} { return []string{} },
			want:     true,
		},
		{
			name:     "non-empty slice",
			getValue: func() interface{} { return []string{"test"} },
			want:     false,
		},
		{
			name:     "zero int",
			getValue: func() interface{} { return 0 },
			want:     true,
		},
		{
			name:     "non-zero int",
			getValue: func() interface{} { return 42 },
			want:     false,
		},
		{
			name:     "false bool",
			getValue: func() interface{} { return false },
			want:     true,
		},
		{
			name:     "true bool",
			getValue: func() interface{} { return true },
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.getValue()
			// Use reflection to get reflect.Value
			rv := reflect.ValueOf(v)
			got := isZeroValue(rv)
			require.Equal(t, tt.want, got)
		})
	}
}
