// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package software

import (
	"fmt"

	"github.com/carabiner-dev/spdx3/base"
	"github.com/carabiner-dev/spdx3/profiles/core"
)

// NewPackage returns a package identified by spdxID.
func NewPackage(spdxID, name string) *Package {
	return &Package{SoftwareArtifact: newSoftwareArtifact(spdxID, PackageType, name)}
}

// NewFile returns a file identified by spdxID.
func NewFile(spdxID, name string) *File {
	return &File{SoftwareArtifact: newSoftwareArtifact(spdxID, FileType, name)}
}

// NewSnippet returns a snippet identified by spdxID, taken from the file it
// names.
func NewSnippet(spdxID, name, fromFile string) *Snippet {
	s := &Snippet{SoftwareArtifact: newSoftwareArtifact(spdxID, SnippetType, name)}
	s.SnippetFromFile = fromFile
	return s
}

// NewSbom returns a bill of materials for software, identified by spdxID.
func NewSbom(spdxID string, sbomType ...SbomType) *Sbom {
	return &Sbom{
		Bom: core.Bom{Bundle: core.Bundle{ElementCollection: core.ElementCollection{
			Node: core.Node{PreNode: base.PreNode{
				SPDXID: spdxID,
				Type:   fmt.Sprintf("%s_%s", Prefix, SbomClass),
			}},
		}}},
		SbomType: sbomType,
	}
}

func newSoftwareArtifact(spdxID, class, name string) SoftwareArtifact {
	return SoftwareArtifact{Artifact: core.Artifact{Node: core.Node{
		PreNode: base.PreNode{
			SPDXID: spdxID,
			Type:   fmt.Sprintf("%s_%s", Prefix, class),
		},
		Name: name,
	}}}
}
