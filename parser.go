// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"io"

	"github.com/carabiner-dev/spdx3/dispatch"
	"github.com/carabiner-dev/spdx3/unmarshal"
)

// ParserOption changes how a parser reads documents.
type ParserOption func(*Parser)

// WithInvalidVocabularyValues keeps values that are not members of their
// vocabulary instead of dropping them, so a document is read exactly as it
// was written. Use it to read documents from a later version of the spec,
// whose vocabularies have entries this one does not know, or to inspect a
// nonconformant document with Validate.
func WithInvalidVocabularyValues() ParserOption {
	return func(p *Parser) { p.keepInvalidVocabularyValues = true }
}

func NewParser(opts ...ParserOption) *Parser {
	unmarshal.SetDefaultDispatcher(dispatch.New())
	p := &Parser{}
	for _, opt := range opts {
		opt(p)
	}
	return p
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
	// keepInvalidVocabularyValues stops the parser dropping property values
	// that are not members of their vocabulary.
	keepInvalidVocabularyValues bool
}

func (p *Parser) Parse(r io.Reader) (*Envelope, error) {
	unmarshal.SetKeepInvalidVocabularyValues(p.keepInvalidVocabularyValues)

	dec := json.NewDecoder(r)
	env := &Envelope{}
	if err := dec.Decode(env); err != nil {
		return nil, err
	}
	return env, nil
}
