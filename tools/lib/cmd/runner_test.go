package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errSimulated = errors.New("simulated failure")

func TestCmdString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cmd  Cmd
		want string
	}{
		{
			name: "simple command",
			cmd:  Cmd{Name: "golangci-lint", Args: []string{"run", "./api/...", "./tools/..."}},
			want: "golangci-lint run ./api/... ./tools/...",
		},
		{
			name: "command with spaces in args",
			cmd:  Cmd{Name: "ifacemaker", Args: []string{"--iface-comment", "This is a comment"}},
			want: "ifacemaker --iface-comment 'This is a comment'",
		},
		{
			name: "command with single quotes in args",
			cmd:  Cmd{Name: "echo", Args: []string{"it's"}},
			want: `echo 'it'\''s'`,
		},
		{
			name: "command with empty arg",
			cmd:  Cmd{Name: "test", Args: []string{"--flag", ""}},
			want: "test --flag ''",
		},
		{
			name: "command with no args",
			cmd:  Cmd{Name: "buf"},
			want: "buf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.cmd.String())
		})
	}
}

func TestDryRunnerRecordsCommands(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := dr.Run(ctx, Cmd{
		Name: "golangci-lint",
		Args: []string{"run", "./api/...", "./tools/..."},
	})
	assert.NoError(t, err)

	err = dr.Run(ctx, Cmd{
		Name: "buf",
		Args: []string{"generate", "--template", "buf.gen.dev.yaml"},
	})
	assert.NoError(t, err)

	require.Len(t, dr.Recorded, 2)
	assert.Equal(t, "golangci-lint", dr.Recorded[0].Name)
	assert.Equal(t, "buf", dr.Recorded[1].Name)
}

func TestDryRunnerStartRecordsAndReturnsProcess(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	proc, err := dr.Start(ctx, Cmd{
		Name: "doppler",
		Args: []string{"run", "--", "go", "run", "./api/cmd/server"},
	})
	require.NoError(t, err)
	assert.Len(t, dr.Recorded, 1)

	// Process.Wait and Process.Stop should be no-ops.
	require.NoError(t, proc.Wait())
	proc.Stop() // should not panic
}

func TestDryRunnerErrorFor(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{
		ErrorFor: map[string]error{
			"golangci-lint": errSimulated,
		},
	}
	ctx := context.Background()

	err := dr.Run(ctx, Cmd{Name: "golangci-lint", Args: []string{"run"}})
	require.ErrorIs(t, err, errSimulated)

	err = dr.Run(ctx, Cmd{Name: "buf", Args: []string{"lint"}})
	require.NoError(t, err)

	assert.Len(t, dr.Recorded, 2)
}

func TestMergeEnv(t *testing.T) {
	t.Parallel()

	parent := []string{"A=1", "B=2", "C=3"}
	extra := []string{"B=override", "D=4"}

	result := mergeEnv(parent, extra)

	assert.Len(t, result, 4)

	resultMap := make(map[string]string, len(result))
	for _, v := range result {
		parts := splitEnvVar(v)
		resultMap[parts[0]] = parts[1]
	}

	assert.Equal(t, "1", resultMap["A"])
	assert.Equal(t, "override", resultMap["B"])
	assert.Equal(t, "3", resultMap["C"])
	assert.Equal(t, "4", resultMap["D"])
}

func splitEnvVar(s string) [2]string {
	i := 0
	for i < len(s) && s[i] != '=' {
		i++
	}

	if i >= len(s) {
		return [2]string{s, ""}
	}

	return [2]string{s[:i], s[i+1:]}
}

func TestReadDopplerVarsDryRunReturnsEmptyMap(t *testing.T) {
	t.Parallel()

	// Use a local DryRunner and set the package-level runner temporarily.
	origRunner := runner
	dr := &DryRunner{}
	runner = dr

	defer func() { runner = origRunner }()

	result := readDopplerVars(context.Background(), dr, "test-project", "dev", "PREFIX_", []string{"DB_HOST"})

	assert.Empty(t, result)
	require.Len(t, dr.Recorded, 1)
	assert.Equal(t, "doppler", dr.Recorded[0].Name)
}
