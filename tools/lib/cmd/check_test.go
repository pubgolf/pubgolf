package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckGo_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := checkGo(context.Background(), dr)
	require.NoError(t, err)

	// checkGo runs golangci-lint cache clean + lint for each of 2 directories.
	require.Len(t, dr.Recorded, 4)

	// First pair: cache clean, then lint ./api/...
	assert.Equal(t, "golangci-lint", dr.Recorded[0].Name)
	assert.Equal(t, []string{"cache", "clean"}, dr.Recorded[0].Args)

	assert.Equal(t, "golangci-lint", dr.Recorded[1].Name)
	assert.Equal(t, []string{"run", "./api/..."}, dr.Recorded[1].Args)

	// Second pair: cache clean, then lint ./tools/...
	assert.Equal(t, "golangci-lint", dr.Recorded[2].Name)
	assert.Equal(t, []string{"cache", "clean"}, dr.Recorded[2].Args)

	assert.Equal(t, "golangci-lint", dr.Recorded[3].Name)
	assert.Equal(t, []string{"run", "./tools/..."}, dr.Recorded[3].Args)
}

func TestCheckGo_StopsOnCacheCleanError(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"golangci-lint": errSimulated,
		},
	}

	err := checkGo(context.Background(), dr)
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)

	// Should have stopped after the first failed command.
	assert.Len(t, dr.Recorded, 1)
}

func TestCheckWeb_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := checkWeb(context.Background(), dr)
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 2)

	assert.Equal(t, "npm", dr.Recorded[0].Name)
	assert.Equal(t, []string{"run", "ci:lint"}, dr.Recorded[0].Args)
	assert.Equal(t, "web-app", dr.Recorded[0].Dir)

	assert.Equal(t, "npm", dr.Recorded[1].Name)
	assert.Equal(t, []string{"run", "check"}, dr.Recorded[1].Args)
	assert.Equal(t, "web-app", dr.Recorded[1].Dir)
}

func TestCheckWeb_StopsOnFirstError(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"npm": errSimulated,
		},
	}

	err := checkWeb(context.Background(), dr)
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)
	assert.Len(t, dr.Recorded, 1)
}

func TestCheckProto_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := checkProto(context.Background(), dr)
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 2)

	assert.Equal(t, "buf", dr.Recorded[0].Name)
	assert.Equal(t, []string{"lint"}, dr.Recorded[0].Args)

	assert.Equal(t, "buf", dr.Recorded[1].Name)
	assert.Equal(t, []string{"format", "./proto/", "--diff", "--exit-code"}, dr.Recorded[1].Args)
}

func TestCheckProto_StopsOnLintError(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"buf": errSimulated,
		},
	}

	err := checkProto(context.Background(), dr)
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)
	assert.Len(t, dr.Recorded, 1)
}
