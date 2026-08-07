// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package spdx3

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carabiner-dev/spdx3/profiles/core"
)

// TestConcurrentParsersDoNotInterfere guards the boundary between parsers.
// Everything a parse needs used to live in package variables that NewParser
// and Parse wrote, so two parsers configured differently would each do the
// other's job: before this was fixed, the strict parser below kept the
// invalid value on roughly 80% of its runs, and the race detector reported
// a data race on every run.
//
// Run with -race to check both properties at once.
func TestConcurrentParsersDoNotInterfere(t *testing.T) {
	strict := NewParser()
	lenient := NewParser(WithInvalidVocabularyValues())

	relationshipType := func(t *testing.T, env *Envelope) core.RelationshipType {
		t.Helper()
		rel, ok := env.Graph[0].(*core.Relationship)
		require.True(t, ok)
		return rel.RelationshipType
	}

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(2)

		go func() {
			defer wg.Done()
			env, err := strict.Parse(strings.NewReader(nonconformantDoc))
			assertNoError(t, err)
			assertEqual(t, core.RelationshipType(""), relationshipType(t, env),
				"the strict parser kept a value the lenient one was asked to keep")
		}()

		go func() {
			defer wg.Done()
			env, err := lenient.Parse(strings.NewReader(nonconformantDoc))
			assertNoError(t, err)
			assertEqual(t, core.RelationshipType("totallyMadeUp"), relationshipType(t, env),
				"the lenient parser dropped a value only the strict one drops")
		}()
	}
	wg.Wait()
}

// require's failure path is not safe to call from a goroutine, since it stops
// the one it runs on rather than the test; these report instead.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("parsing: %v", err)
	}
}

func assertEqual[T comparable](t *testing.T, want, got T, message string) {
	t.Helper()
	if want != got {
		t.Errorf("%s: wanted %v, got %v", message, want, got)
	}
}
