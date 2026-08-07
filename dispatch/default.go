// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package dispatch

import (
	"maps"

	"github.com/carabiner-dev/spdx3/profiles/ai"
	"github.com/carabiner-dev/spdx3/profiles/build"
	"github.com/carabiner-dev/spdx3/profiles/core"
	"github.com/carabiner-dev/spdx3/profiles/dataset"
	"github.com/carabiner-dev/spdx3/profiles/expandedlicensing"
	"github.com/carabiner-dev/spdx3/profiles/extension"
	"github.com/carabiner-dev/spdx3/profiles/security"
	"github.com/carabiner-dev/spdx3/profiles/simplelicensing"
	"github.com/carabiner-dev/spdx3/profiles/software"
	"github.com/carabiner-dev/spdx3/types"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

// Classes returns the registry mapping every SPDX type name this library
// knows to a prototype of the Go type modelling it, in both its bare and
// profile-prefixed spellings. It is plain data, so an unmarshaler can hold
// it without this package's behaviour and without an import cycle.
func Classes() map[string]types.Node {
	classes := map[string]types.Node{}

	// Add the types for all profiles
	maps.Copy(classes, core.Profile.Classes)
	maps.Copy(classes, ai.Profile.Classes)
	maps.Copy(classes, build.Profile.Classes)
	maps.Copy(classes, dataset.Profile.Classes)
	maps.Copy(classes, expandedlicensing.Profile.Classes)
	maps.Copy(classes, extension.Profile.Classes)
	maps.Copy(classes, security.Profile.Classes)
	maps.Copy(classes, simplelicensing.Profile.Classes)
	maps.Copy(classes, software.Profile.Classes)

	return classes
}

// New returns a dispatcher over every class this library knows.
func New() types.Dispatcher {
	return unmarshal.New(Classes(), unmarshal.Options{})
}
