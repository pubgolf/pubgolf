package cmd

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceSharedDBPort_ReplacesPort(t *testing.T) {
	t.Parallel()

	env := []string{
		"SOME_VAR=hello",
		"PUBGOLF_SHARED_DB_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable",
		"OTHER_VAR=world",
	}

	result := replaceSharedDBPort(env, 5)

	// The port should be 5432 + 5 = 5437.
	assert.Contains(t, result[1], "5437")
	assert.Contains(t, result[1], "PUBGOLF_SHARED_DB_URL=")
	// Other vars should be unchanged.
	assert.Equal(t, "SOME_VAR=hello", result[0])
	assert.Equal(t, "OTHER_VAR=world", result[2])
}

func TestReplaceSharedDBPort_NoMatch(t *testing.T) {
	t.Parallel()

	env := []string{"SOME_VAR=hello", "OTHER_VAR=world"}
	result := replaceSharedDBPort(env, 5)

	assert.Equal(t, env, result)
}

func TestReplaceSharedDBPort_InvalidURL(t *testing.T) {
	t.Parallel()

	env := []string{"PUBGOLF_SHARED_DB_URL=://invalid"}
	result := replaceSharedDBPort(env, 5)

	// Should return env unchanged when URL can't be parsed.
	assert.Equal(t, env, result)
}

func TestPreflight_SucceedsWhenPortOpen(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Start a listener on a random port to simulate a running DB.
	var lc net.ListenConfig

	ln, err := lc.Listen(ctx, "tcp", "localhost:0")
	require.NoError(t, err)

	defer ln.Close()

	// Extract the port from the listener's address.
	addr, ok := ln.Addr().(*net.TCPAddr)
	require.True(t, ok, "expected *net.TCPAddr")

	offset := addr.Port - 5432

	err = preflight(ctx, offset)
	assert.NoError(t, err)
}

func TestPreflight_FailsWhenPortClosed(t *testing.T) {
	t.Parallel()

	// Use a port that is almost certainly not listening.
	// Offset of 499 → port 5931 — unlikely to be in use in a test environment.
	err := preflight(context.Background(), 499)
	require.Error(t, err)
	require.ErrorIs(t, err, errDBNotRunning)
}
