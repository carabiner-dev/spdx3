package spdx3

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrUnsupportedNodeType = errors.New("unsupported node type")

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
