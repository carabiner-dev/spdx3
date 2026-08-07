# Carabienr SPDX3 Writer/Parser

A Go library for reading and writing [SPDX 3](https://spdx.github.io/spdx-spec/)
documents.

SPDX 3 documents are JSON-LD, which makes them more harder to handle than plain
JSON: the same element can appear as a full object in one place and as a bare
IRI reference in another, identifiers come in several flavours, and properties
outside the Core profile carry a profile prefix. This library hides that,
parsing a document into typed Go structs and writing it back as a document the
SPDX tools accept.

This library targets **SPDX 3.0.1** at the moment. At least at the time of the
initial release, every one of the official SPDX example documents round trips
through it without losing or changing anything. We also check the output
on every pull request against the official JSON schema and the SPDX Java
tools. See [Conformance](#conformance).

## Install

```sh
go get github.com/carabiner-dev/spdx3
```

## Reading a Document

`Parse` returns an `Envelope`: the document's `@context` and its graph of nodes,
each already typed to its SPDX class.

```go
file, err := os.Open("document.spdx.json")
if err != nil {
    return err
}
defer file.Close()

env, err := spdx3.NewParser().Parse(file)
if err != nil {
    return err
}

fmt.Println("spec version:", env.Context.Version())

for _, node := range env.Graph {
    switch n := node.(type) {
    case *software.Package:
        fmt.Printf("package %s %s\n", n.GetName(), n.PackageVersion)
        for _, im := range n.VerifiedUsing {
            if h, ok := im.(*core.Hash); ok {
                fmt.Printf("  %s: %s\n", h.Algorithm, h.HashValue)
            }
        }
        // A reference resolves to the element itself, so following one
        // reads its properties without a lookup of your own.
        if ci, ok := n.GetCreationInfo().(*core.CreationInfo); ok {
            fmt.Printf("  created %s by %s\n", ci.Created, ci.CreatedBy[0].GetName())
        }
    case *core.Relationship:
        fmt.Printf("%s %s %d element(s)\n",
            n.From.GetSPDXID(), n.RelationshipType, len(n.To))
    }
}
```

Every node satisfies `types.Node`, so `GetSPDXID`, `GetID`, `GetType`, `GetName`
and `GetCreationInfo` work without a type switch when that is all you need.

## Writing a Document

Build the elements, say how they relate, and render. A document needs a shape:
an `SpdxDocument` names what it is about through its root elements, and
relationships say how the rest connect.

```go
alice := core.NewPerson("https://example.com/spdx/alice", "Alice")
tool := core.NewTool("https://example.com/spdx/tool", "my-generator")

pkg := software.NewPackage("https://example.com/spdx/pkg", "example-lib")
pkg.PackageVersion = "1.4.2"
pkg.DownloadLocation = "https://example.com/example-lib-1.4.2.tar.gz"
pkg.PrimaryPurpose = software.SoftwarePurposeLibrary
pkg.VerifiedUsing = []core.IntegrityMethodDescendant{
    core.NewHash(core.HashAlgorithmSha256,
        "5d41402abc4b2a76b9719d911017c592f0e8b1b45c0f47b09fb8f0e2e0d9c0aa"),
}

file := software.NewFile("https://example.com/spdx/file", "./src/main.go")

doc := core.NewSpdxDocument("https://example.com/spdx/document")
doc.AddRootElement(pkg)
doc.ProfileConformance = []core.ProfileIdentifierType{
    core.ProfileIdentifierTypeCore, core.ProfileIdentifierTypeSoftware,
}

env := spdx3.NewEnvelope()
env.Graph.AddNode(alice, tool, doc, pkg, file)

env.Graph.Relate("https://example.com/spdx/describes",
    doc, core.RelationshipTypeDescribes, pkg)
env.Graph.Relate("https://example.com/spdx/contains",
    pkg, core.RelationshipTypeContains, file)

// Said once, and shared by every element that has none of its own.
creation := core.NewCreationInfo(time.Now(), alice)
creation.CreatedUsing = []types.Node{tool}
env.Graph.SetCreationInfo(creation)

if err := (&spdx3.Renderer{}).Render(env, os.Stdout); err != nil {
    return err
}
```

The constructors set each element's type discriminator, which a document is
read by and which a struct literal has to get right by hand. Every field is
still reachable: they return the element, not a builder.

An element referenced from a property is written as a reference when the graph
carries it, so each one appears once, in the graph, however many times it is
named.

Timestamps go through `types.DateTime`, which writes the lexical form SPDX
requires — whole seconds, in UTC — whatever precision and zone the `time.Time`
you hand it carries.

## What is covered

All nine profiles that carry classes in 3.0.1, with 84 type names registered for
dispatch, in both their bare and profile-prefixed spellings:

| Profile | Package |
| --- | --- |
| Core | `profiles/core` |
| Software | `profiles/software` |
| Licensing (simple and expanded) | `profiles/simplelicensing`, `profiles/expandedlicensing` |
| Security | `profiles/security` |
| AI | `profiles/ai` |
| Dataset | `profiles/dataset` |
| Build | `profiles/build` |
| Extension | `profiles/extension` |

Vocabularies are exposed as typed string constants (`core.RelationshipType`,
`core.HashAlgorithm`, `software.SoftwarePurpose`, …) and match the 3.0.1 model
exactly. The individuals the spec predefines are available as values with their
IRIs: `core.NoneElement`, `core.NoAssertionElement`, `core.SpdxOrganization`,
`expandedlicensing.NoneLicense` and `expandedlicensing.NoAssertionLicense`.

Beyond vocabularies, nothing is validated: missing required properties and
wrong cardinalities are read and written as they come.

## Vocabulary Values

A property drawn from a vocabulary only admits the values that vocabulary
lists. Parsing drops anything else, so a value that reaches your code is one
the specification defines:

```go
// "relationshipType": "totallyMadeUp" in the document
rel.RelationshipType == ""   // dropped
```

Every vocabulary type answers `IsValid()`, and the list behind it answers
`Contains`, so you can check a value before assigning it.

Dropping loses whatever the document said, which matters when the document
comes from a later version of the spec: SPDX 3.1 adds seventeen relationship
types, and to a 3.0.1 vocabulary those look like values to discard. Parse with
`WithInvalidVocabularyValues` to keep them, and use `Validate` to see what is
not recognized:

```go
env, err := spdx3.NewParser(spdx3.WithInvalidVocabularyValues()).Parse(file)
if err != nil {
    return err
}

for _, finding := range spdx3.Validate(env) {
    fmt.Println(finding)
    // https://example.com/rel (Relationship): relationshipType:
    //   "totallyMadeUp" is not a member of the vocabulary this property draws from
}
```

With the values kept the document also round-trips as written, rather than
coming back out without them.

## Notes on the JSON-LD shapes

- **References and inline objects:** A field holding another node accepts both a
  bare IRI and a full object, and both end up as the node itself: a reference
  to something the document carries is resolved to that node, so following a
  property leads to the element rather than to its identifier, and two
  properties naming the same element get the same value. A reference the
  document has no node for (one naming an element in another document, or an
  individual the spec predefines) stays a `types.NodeRef`, which carries only
  the identifier.
- **Blank nodes:** `_:` labels are kept verbatim, so a reference still matches
  the node it points at after a round trip.
- **Contexts:** `@context` may be a URL, an array or an inline object. All three
  parse, and the array and object forms are preserved exactly. `Context.Version`
  reports the spec version the context pins.
- **Documents without `@graph`:** A lone element as the document root is read
  into a single-node graph, and written back inside a `@graph`.
- **`@type`:** The context aliases `type` to `@type`, so a class named either
  way is read, what gets written is always the plain spelling.
- **Properties we do not know** are ignored, so a document written against a
  later version of the specification still reads. A node where *nothing*
  binds is reported instead, since reading it would quietly yield an empty
  element — that is what a document carrying fully expanded JSON-LD IRIs,
  rather than the names the SPDX context defines, looks like. Expanding such
  a document is out of scope; compact it first.

## Conformance

The tests check two different things. `go test ./...` parses, renders and
reparses every document under `testdata/`, and asserts the result says exactly
what the input said: every node still found by its identifier, no value changed,
no property gained or lost, and every reference still resolving. That runs over
the five documents in `testdata/` and the 26 SPDX example documents vendored
in `testdata/corpus`.

`hack/verify-spdx.sh` goes further and checks the documents this library *writes*
against the SPDX project's own tooling: the official 3.0.1 JSON schema and the
`Verify` command of the [SPDX Java tools](https://github.com/spdx/tools-java).
Besides the rendered corpus it checks a document built out of Go values, since
re-rendering a document that was read from a file cannot exercise the code that
turns Go values into SPDX ones. It needs `java` and `check-jsonschema`:

```sh
./hack/verify-spdx.sh
```

For the rendered documents both checks also compare against the input, and only
report a failure where the tools accept the input but reject our output, so a
defect in an upstream example cannot fail the build. The authored document has
no input to compare against and simply has to be valid.

To refresh the vendored examples after the SPDX project publishes new ones:

```sh
./hack/update-spdx-examples.sh
```

The examples are the SPDX project's, provided under CC0-1.0.

## SPDX 3.1

3.1 will be additive over 3.0.1: no class, property, vocabulary entry or individual
was removed or renamed, so the types here stay valid. A 3.1 document whose
classes exist in 3.0.1 already parses, since dispatch keys on the type name
rather than the spec version, and `Context.Version` recognizes the 3.1 context
URLs. Profiles from 3.1 are not implemented yet.

## License

Built with ♥️ by Carabiner Systems, Inc and released under the Apache-2.0 license.
Patches and issues are welcome!

