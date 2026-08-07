// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"encoding/json"
	"fmt"
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

// Unmarshaling a node on its own, with encoding/json rather than a Parser,
// goes through the package-level dispatcher. Set it once here; a Parser
// carries its own and never writes this.
func init() {
	unmarshal.SetDefaultDispatcher(dispatch.New())
}

func NewParser(opts ...ParserOption) *Parser {
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
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading the document: %w", err)
	}

	// The unmarshaler carries everything this parse needs, so parsers
	// configured differently can run at the same time.
	nu := unmarshal.New(dispatch.Classes(), unmarshal.Options{
		KeepInvalidVocabularyValues: p.keepInvalidVocabularyValues,
	})

	env := &Envelope{}
	if err := env.unmarshalWith(data, nu); err != nil {
		return nil, err
	}
	return env, nil
}
