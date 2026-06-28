package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

// This file ensures test dependencies are tracked in go.mod.
// It will be replaced by proper tests in later tasks.

func TestPlaceholder(t *testing.T) {
	assert.True(t, true)
	rapid.Check(t, func(t *rapid.T) {
		n := rapid.IntRange(0, 100).Draw(t, "n")
		assert.GreaterOrEqual(t, n, 0)
	})
}
