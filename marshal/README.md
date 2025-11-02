# SPDX3 Marshal Package

The `marshal` package provides functionality to serialize SPDX3 nodes and graphs to JSON format. It's the counterpart to the `unmarshal` package and handles the reverse operation.

## Overview

When serializing an `spdx.Graph` to JSON:
- **Top-level nodes** in the graph are serialized fully (complete struct to JSON)
- **Nested nodes** or descendants within those top-level nodes are serialized as just their ID strings (SPDXID references)

This approach maintains the JSON-LD structure where nested nodes are referenced by their IDs rather than embedded as complete objects.

## Key Features

- **Automatic nested node detection**: Uses reflection to detect fields that implement `types.Node` interface
- **Reference serialization**: Nested nodes are automatically serialized as ID strings
- **Embedded struct handling**: Properly handles embedded structs like `PreNode`, flattening their fields into the parent JSON object
- **Omitempty support**: Respects `json:",omitempty"` tags to exclude zero values
- **Type-safe**: Works with any type that implements the `types.Node` interface

## Usage

### Basic Usage with Graph

The `Graph` type has a built-in `MarshalJSON` method that automatically uses the marshal package:

```go
import (
    "encoding/json"
    "github.com/carabiner-dev/spdx3"
)

// Create a graph with nodes
graph := spdx3.Graph{
    creationInfo,
    person,
    organization,
}

// Marshal to JSON
data, err := json.Marshal(graph)
if err != nil {
    log.Fatal(err)
}

// Or marshal the entire Envelope
envelope := &spdx3.Envelope{
    Context: "https://spdx.org/rdf/3.0.1/spdx-context.jsonld",
    Graph:   graph,
}

data, err := json.MarshalIndent(envelope, "", "  ")
```

### Manual Usage with NodeMarshaler

You can also use the `NodeMarshaler` directly:

```go
import "github.com/carabiner-dev/spdx3/marshal"

marshaler := &marshal.NodeMarshaler{}
data, err := marshaler.MarshalNode(node)
```

## How It Works

### Example Input

Consider a `CreationInfo` node with a nested `Person` in the `CreatedBy` field:

```go
person := &core.Person{
    Agent: core.Agent{
        Node: core.Node{
            PreNode: base.PreNode{
                ID:     "https://spdx.org/spdxdocs/Person1-abc123",
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
    },
    SpecVersion: "3.0.1",
    CreatedBy:   []core.AgentDescendant{person},
    Created:     &now,
}
```

### Example Output

When marshaled as part of a graph, the nested `Person` in `CreatedBy` is serialized as just its ID:

```json
{
  "@id": "_:creationinfo",
  "type": "CreationInfo",
  "specVersion": "3.0.1",
  "createdBy": [
    "https://spdx.org/spdxdocs/Person1-abc123"
  ],
  "created": "2024-05-31T00:00:00Z"
}
```

The full `Person` object appears separately in the graph:

```json
{
  "@id": "https://spdx.org/spdxdocs/Person1-abc123",
  "type": "Person",
  "spdxID": "SPDXRef-Person1",
  "name": "John Doe"
}
```

## Nested Node Detection

The marshaler automatically detects and handles:

1. **Single Node fields**: `CreationInfo *CreationInfo` → serialized as `"_:creationinfo"`
2. **Node slices**: `CreatedBy []AgentDescendant` → serialized as `["id1", "id2"]`
3. **NodeRef types**: Already references, serialized as their ID string
4. **Embedded structs**: `base.PreNode` fields are flattened into the parent object

## Roundtrip Support

The marshal package is designed to support perfect roundtrips:

```go
// Parse JSON
env1, err := parser.Parse(reader)

// Marshal back to JSON
data, err := json.Marshal(env1)

// Parse again
env2, err := parser.Parse(bytes.NewReader(data))

// env1 and env2 should be equivalent
```

See the `TestRoundtrip` test in `parser_test.go` for a complete example.

## Implementation Details

The marshaler uses Go's `reflect` package to:
1. Iterate through struct fields
2. Check JSON tags for field names and options
3. Detect if fields implement the `types.Node` interface
4. Recursively handle embedded structs
5. Serialize Node types as ID references

Key functions:
- `NodeMarshaler.MarshalNode()`: Entry point for marshaling any node
- `marshalToMap()`: Converts a struct to a map using reflection
- `marshalSingleNode()`: Handles single node fields (converts to ID string)
- `marshalNodeSlice()`: Handles slices of nodes (converts to array of ID strings)
- `isZeroValue()`: Checks if a value should be omitted with `omitempty`

## Testing

The package includes comprehensive tests:
- `marshal/node_test.go`: Tests for the NodeMarshaler
- `parser_test.go`: Tests for Graph marshaling and roundtrip functionality

Run tests with:
```bash
go test ./marshal/...
go test -run "TestGraphMarshalJSON|TestRoundtrip"
```
