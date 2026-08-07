// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type colour string

func TestVocabularyContains(t *testing.T) {
	colours := Vocabulary[colour]{"red", "green"}

	require.True(t, colours.Contains("red"))
	require.True(t, colours.Contains("green"))
	require.False(t, colours.Contains("mauve"))
	require.False(t, colours.Contains(""), "the empty value is not a member")
	require.False(t, Vocabulary[colour]{}.Contains("red"))
}
