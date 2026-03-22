package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateProto_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := generateProto(context.Background(), dr)
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 1)

	cmd := dr.Recorded[0]
	assert.Equal(t, "buf", cmd.Name)
	assert.Equal(t, []string{"generate", "--template", filepath.FromSlash("buf.gen.dev.yaml")}, cmd.Args)
}

func TestGenerateProto_Error(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"buf": errSimulated,
		},
	}

	err := generateProto(context.Background(), dr)
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)
}

func TestGenerateSQLc_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := generateSQLc(context.Background(), dr)
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 1)

	cmd := dr.Recorded[0]
	assert.Equal(t, "sqlc", cmd.Name)
	assert.Equal(t, []string{"generate", "--file", filepath.FromSlash("api/internal/db/sqlc.yaml")}, cmd.Args)
}

func TestGenerateEnum_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := generateEnum(context.Background(), dr, "ScoringCategory", filepath.FromSlash("./api/internal/lib/models"))
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 1)

	cmd := dr.Recorded[0]
	assert.Equal(t, "enumer", cmd.Name)
	assert.Equal(t, []string{"-sql", "-transform", "snake-upper", "-type", "ScoringCategory", filepath.FromSlash("./api/internal/lib/models")}, cmd.Args)
}

func TestGenerateEnum_Error(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"enumer": errSimulated,
		},
	}

	err := generateEnum(context.Background(), dr, "ScoringCategory", "./api/internal/lib/models")
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)
}

func TestGenerateMock_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	err := generateMock(context.Background(), dr, "Querier", filepath.FromSlash("api/internal/lib/dao/internal/dbc/"))
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 1)

	cmd := dr.Recorded[0]
	assert.Equal(t, "mockery", cmd.Name)
	assert.Contains(t, cmd.Args, "--dir")
	assert.Contains(t, cmd.Args, filepath.FromSlash("api/internal/lib/dao/internal/dbc/"))
	assert.Contains(t, cmd.Args, "--name")
	assert.Contains(t, cmd.Args, "Querier")
	assert.Contains(t, cmd.Args, "--filename")
	assert.Contains(t, cmd.Args, "gen_mock.go")
	assert.Contains(t, cmd.Args, "--inpackage")
}

func TestGenerateMock_Error(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"mockery": errSimulated,
		},
	}

	err := generateMock(context.Background(), dr, "Querier", "api/internal/lib/dao/internal/dbc/")
	require.Error(t, err)
	require.ErrorIs(t, err, errSimulated)
}
