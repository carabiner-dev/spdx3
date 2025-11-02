// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"io"

	"github.com/carabiner-dev/spdx3/dispatch"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

func NewParser() *Parser {
	unmarshal.SetDefaultDispatcher(dispatch.New())
	return &Parser{}
}

type Profile interface {
	Prefix() string
}

type Renderer struct{}

func (r *Renderer) Render(env *Envelope, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(*env)
}

type Parser struct {
}

func (p *Parser) Parse(r io.Reader) (*Envelope, error) {
	dec := json.NewDecoder(r)
	env := &Envelope{}
	if err := dec.Decode(env); err != nil {
		return nil, err
	}
	return env, nil
}
