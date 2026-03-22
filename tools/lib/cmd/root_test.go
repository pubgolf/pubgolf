package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveProjectRoot_FromRoot(t *testing.T) { //nolint:paralleltest // t.Chdir is incompatible with t.Parallel.
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o600))

	t.Chdir(dir)

	root, err := resolveProjectRoot()
	require.NoError(t, err)
	assert.Equal(t, dir, root)
}

func TestResolveProjectRoot_FromSubdir(t *testing.T) { //nolint:paralleltest // t.Chdir is incompatible with t.Parallel.
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o600))

	subdir := filepath.Join(dir, "a", "b", "c")
	require.NoError(t, os.MkdirAll(subdir, 0o755))

	t.Chdir(subdir)

	root, err := resolveProjectRoot()
	require.NoError(t, err)
	assert.Equal(t, dir, root)
}

func TestResolveProjectRoot_NotFound(t *testing.T) { //nolint:paralleltest // t.Chdir is incompatible with t.Parallel.
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	t.Chdir(dir)

	_, err = resolveProjectRoot()
	require.ErrorIs(t, err, errProjectRootNotFound)
}
