package cmd

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerRun_CommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ep := &dryRunProvider{}

	// Save and restore package-level projectRoot since dockerRun uses it.
	origRoot := projectRoot
	projectRoot = "/tmp/test-project"

	defer func() { projectRoot = origRoot }()

	err := dockerRun(context.Background(), dr, ep, "pubgolf-api-server", "dev", "api-db")
	require.NoError(t, err)

	require.Len(t, dr.Recorded, 1)

	cmd := dr.Recorded[0]
	assert.Equal(t, "docker", cmd.Name)
	assert.Contains(t, cmd.Args, "compose")
	assert.Contains(t, cmd.Args, "--file")
	assert.Contains(t, cmd.Args, filepath.FromSlash("./infra/docker-compose.dev.yaml"))
	assert.Contains(t, cmd.Args, "--project-name")
	assert.Contains(t, cmd.Args, "up")
	assert.Contains(t, cmd.Args, "--detach")
	assert.Contains(t, cmd.Args, "api-db")
}

func TestDockerRun_InjectsPortEnv(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ep := &dryRunProvider{}

	origRoot := projectRoot
	projectRoot = "/tmp/test-project"

	defer func() { projectRoot = origRoot }()

	err := dockerRun(context.Background(), dr, ep, "pubgolf-api-server", "dev", "api-db")
	require.NoError(t, err)
	require.Len(t, dr.Recorded, 1)

	// Verify port-related env vars are injected.
	env := dr.Recorded[0].Env
	hasDBPort := false
	hasAPIPort := false
	hasDataPath := false

	for _, e := range env {
		if len(e) >= len("PUBGOLF_DB_PORT=") && e[:len("PUBGOLF_DB_PORT=")] == "PUBGOLF_DB_PORT=" {
			hasDBPort = true
		}

		if len(e) >= len("PUBGOLF_PORT=") && e[:len("PUBGOLF_PORT=")] == "PUBGOLF_PORT=" {
			hasAPIPort = true
		}

		if len(e) >= len("PUBGOLF_DB_HOST_DATA_PATH=") && e[:len("PUBGOLF_DB_HOST_DATA_PATH=")] == "PUBGOLF_DB_HOST_DATA_PATH=" {
			hasDataPath = true
		}
	}

	assert.True(t, hasDBPort, "expected PUBGOLF_DB_PORT in env")
	assert.True(t, hasAPIPort, "expected PUBGOLF_PORT in env")
	assert.True(t, hasDataPath, "expected PUBGOLF_DB_HOST_DATA_PATH in env")
}

func TestDockerRun_EnvProviderError(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ep := &DryRunner{
		ErrorFor: map[string]error{
			"doppler": errSimulated,
		},
	}

	// Use a DopplerProvider that will fail since DryRunner can't capture stdout.
	dp := &DopplerProvider{Runner: ep}

	origRoot := projectRoot
	projectRoot = "/tmp/test-project"

	defer func() { projectRoot = origRoot }()

	err := dockerRun(context.Background(), dr, dp, "pubgolf-api-server", "dev", "api-db")
	require.Error(t, err)
}
