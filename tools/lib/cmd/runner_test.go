package cmd

import (
	"context"
	"errors"
	"slices"
	"testing"
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

			got := tt.cmd.String()
			if got != tt.want {
				t.Errorf("Cmd.String() = %q, want %q", got, tt.want)
			}
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
	if err != nil {
		t.Fatalf("DryRunner.Run() returned error: %v", err)
	}

	err = dr.Run(ctx, Cmd{
		Name: "buf",
		Args: []string{"generate", "--template", "buf.gen.dev.yaml"},
	})
	if err != nil {
		t.Fatalf("DryRunner.Run() returned error: %v", err)
	}

	if len(dr.Recorded) != 2 {
		t.Fatalf("expected 2 recorded commands, got %d", len(dr.Recorded))
	}

	if dr.Recorded[0].Name != "golangci-lint" {
		t.Errorf("first command name = %q, want %q", dr.Recorded[0].Name, "golangci-lint")
	}

	if dr.Recorded[1].Name != "buf" {
		t.Errorf("second command name = %q, want %q", dr.Recorded[1].Name, "buf")
	}
}

func TestDryRunnerStartRecordsAndReturnsProcess(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	proc, err := dr.Start(ctx, Cmd{
		Name: "doppler",
		Args: []string{"run", "--", "go", "run", "./api/cmd/server"},
	})
	if err != nil {
		t.Fatalf("DryRunner.Start() returned error: %v", err)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	// Process.Wait and Process.Stop should be no-ops.
	waitErr := proc.Wait()
	if waitErr != nil {
		t.Errorf("dryProcess.Wait() returned error: %v", waitErr)
	}

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
	if !errors.Is(err, errSimulated) {
		t.Errorf("expected simulated error, got: %v", err)
	}

	err = dr.Run(ctx, Cmd{Name: "buf", Args: []string{"lint"}})
	if err != nil {
		t.Errorf("expected nil error for buf, got: %v", err)
	}

	if len(dr.Recorded) != 2 {
		t.Errorf("expected 2 recorded commands, got %d", len(dr.Recorded))
	}
}

func TestCheckGoCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := checkGo(ctx, dr)
	if err != nil {
		t.Fatalf("checkGo() returned error: %v", err)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	cmd := dr.Recorded[0]
	if cmd.Name != "golangci-lint" {
		t.Errorf("command name = %q, want %q", cmd.Name, "golangci-lint")
	}

	if len(cmd.Args) < 1 || cmd.Args[0] != "run" {
		t.Errorf("expected first arg to be 'run', got %v", cmd.Args)
	}
}

func TestCheckProtoCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := checkProto(ctx, dr)
	if err != nil {
		t.Fatalf("checkProto() returned error: %v", err)
	}

	if len(dr.Recorded) != 2 {
		t.Fatalf("expected 2 recorded commands, got %d", len(dr.Recorded))
	}

	if dr.Recorded[0].Name != "buf" || dr.Recorded[0].Args[0] != "lint" {
		t.Errorf("first command should be 'buf lint', got %s %v", dr.Recorded[0].Name, dr.Recorded[0].Args)
	}

	if dr.Recorded[1].Name != "buf" || dr.Recorded[1].Args[0] != "format" {
		t.Errorf("second command should be 'buf format', got %s %v", dr.Recorded[1].Name, dr.Recorded[1].Args)
	}
}

func TestCheckWebCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := checkWeb(ctx, dr)
	if err != nil {
		t.Fatalf("checkWeb() returned error: %v", err)
	}

	if len(dr.Recorded) != 2 {
		t.Fatalf("expected 2 recorded commands, got %d", len(dr.Recorded))
	}

	if dr.Recorded[0].Dir != "web-app" {
		t.Errorf("first command dir = %q, want %q", dr.Recorded[0].Dir, "web-app")
	}

	if dr.Recorded[1].Dir != "web-app" {
		t.Errorf("second command dir = %q, want %q", dr.Recorded[1].Dir, "web-app")
	}
}

func TestGenerateProtoCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := generateProto(ctx, dr)
	if err != nil {
		t.Fatalf("generateProto() returned error: %v", err)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	cmd := dr.Recorded[0]
	if cmd.Name != "buf" {
		t.Errorf("command name = %q, want %q", cmd.Name, "buf")
	}
}

func TestGenerateSQLcCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := generateSQLc(ctx, dr)
	if err != nil {
		t.Fatalf("generateSQLc() returned error: %v", err)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	if dr.Recorded[0].Name != "sqlc" {
		t.Errorf("command name = %q, want %q", dr.Recorded[0].Name, "sqlc")
	}
}

func TestGenerateEnumCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}
	ctx := context.Background()

	err := generateEnum(ctx, dr, "ScoringCategory", "./api/internal/lib/models")
	if err != nil {
		t.Fatalf("generateEnum() returned error: %v", err)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	cmd := dr.Recorded[0]
	if cmd.Name != "enumer" {
		t.Errorf("command name = %q, want %q", cmd.Name, "enumer")
	}
}

func TestMergeEnv(t *testing.T) {
	t.Parallel()

	parent := []string{"A=1", "B=2", "C=3"}
	extra := []string{"B=override", "D=4"}

	result := mergeEnv(parent, extra)

	expected := map[string]string{
		"A": "1",
		"B": "override",
		"C": "3",
		"D": "4",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d env vars, got %d", len(expected), len(result))
	}

	for _, v := range result {
		parts := splitEnvVar(v)
		if want, ok := expected[parts[0]]; ok {
			if parts[1] != want {
				t.Errorf("env %s = %q, want %q", parts[0], parts[1], want)
			}
		} else {
			t.Errorf("unexpected env var: %s", v)
		}
	}
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

func TestDopplerDockerStopCommandConstruction(t *testing.T) {
	t.Parallel()

	dr := &DryRunner{}

	dopplerDockerStop(context.Background(), dr, "test-project", "dev")

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	cmd := dr.Recorded[0]
	if cmd.Name != "doppler" {
		t.Errorf("command name = %q, want %q", cmd.Name, "doppler")
	}

	if !slices.Contains(cmd.Args, "down") {
		t.Error("expected 'down' in command args")
	}
}

func TestReadDopplerVarsDryRunReturnsEmptyMap(t *testing.T) {
	t.Parallel()

	// Use a local DryRunner and set the package-level runner temporarily.
	origRunner := runner
	dr := &DryRunner{}
	runner = dr

	defer func() { runner = origRunner }()

	result := readDopplerVars(context.Background(), dr, "test-project", "dev", "PREFIX_", []string{"DB_HOST"})

	if len(result) != 0 {
		t.Errorf("expected empty map in dry-run mode, got %v", result)
	}

	if len(dr.Recorded) != 1 {
		t.Fatalf("expected 1 recorded command, got %d", len(dr.Recorded))
	}

	if dr.Recorded[0].Name != "doppler" {
		t.Errorf("expected doppler command to be recorded, got %q", dr.Recorded[0].Name)
	}
}
