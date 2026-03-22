# Plan 09: Test Suite + CLAUDE.md Updates

**Depends on:** 01-runner-interface, 02-worktree-identity, 05-error-handling
**Branch:** `devctrl/09-tests-and-claudemd`
**PR scope:** `tools/lib/cmd/*_test.go`, `CLAUDE.md`

## Objective

Comprehensive test suite for devctrl using the DryRunner, plus CLAUDE.md
updates for the new commands and agent guidance.

## Prerequisites

Plans 01 (Runner/DryRunner), 02 (worktree identity), and 05 (error handling)
must be merged. Tests use DryRunner for command recording, worktree functions
for isolation tests, and error classification for exit code tests.

## Steps

### 1. Command Construction Tests — `runner_test.go`

Test each command handler's subprocess invocation. Use DryRunner to capture
what would be executed and assert on name, args, dir, env.

```go
func TestCheckGo_CommandConstruction(t *testing.T) {
    dr := &DryRunner{}
    err := checkGo(context.Background(), dr)
    require.NoError(t, err)
    require.Len(t, dr.Recorded, 1)
    cmd := dr.Recorded[0]
    assert.Equal(t, "golangci-lint", cmd.Name)
    assert.Contains(t, cmd.Args, "--exclude-dirs")
    assert.Contains(t, cmd.Args, `\.worktrees`)
}

func TestCheckWeb_CommandConstruction(t *testing.T) {
    dr := &DryRunner{}
    err := checkWeb(context.Background(), dr)
    require.NoError(t, err)
    require.Len(t, dr.Recorded, 2)  // ci:lint + check
    assert.Equal(t, "npm", dr.Recorded[0].Name)
    assert.Equal(t, "web-app", dr.Recorded[0].Dir)
}

func TestCheckProto_CommandConstruction(t *testing.T) {
    dr := &DryRunner{}
    err := checkProto(context.Background(), dr)
    require.NoError(t, err)
    require.Len(t, dr.Recorded, 2)  // buf lint + buf format
}
```

### 2. Generate Command Tests — `generate_test.go`

```go
func TestGenerateProto_CommandConstruction(t *testing.T) {
    dr := &DryRunner{}
    err := generateProto(context.Background(), dr)
    require.NoError(t, err)
    require.Len(t, dr.Recorded, 1)
    assert.Equal(t, "buf", dr.Recorded[0].Name)
    assert.Contains(t, dr.Recorded[0].Args, "buf.gen.dev.yaml")
}

func TestGenerateSQLc_CommandConstruction(t *testing.T) { ... }
func TestGenerateEnum_CommandConstruction(t *testing.T) { ... }
func TestGenerateMock_SequentialAfterDBC(t *testing.T) { ... }
```

### 3. Flag Propagation Tests — `test_cmd_test.go`

Verify that CLI flags modify the constructed command correctly:

```go
func TestTestCmd_VerboseFlag(t *testing.T) {
    // Set up cobra command with --verbose=true
    // Invoke the test handler
    // Assert "-v" appears in the recorded command args
}

func TestTestCmd_CoverageFlag(t *testing.T) { ... }
func TestTestCmd_LocalFlag(t *testing.T) { ... }
```

### 4. Worktree Isolation Tests — `worktree_integration_test.go`

These tests verify that worktree identity is correctly threaded into command
construction. They mock the worktree slug and verify env injection:

```go
func TestDockerRun_WorktreeProjectName(t *testing.T) {
    // Set worktree slug to "fix-auth"
    // Invoke dockerRun with DryRunner
    // Assert "--project-name" "pubgolf-fix-auth" in recorded args
}

func TestTestCmd_InjectsWorktreeSlug(t *testing.T) {
    // Set worktree slug to "fix-auth"
    // Invoke test handler with DryRunner
    // Assert "PUBGOLF_WORKTREE_SLUG=fix-auth" in recorded env
}

func TestTestCmd_InjectsOffsetPort(t *testing.T) {
    // Set worktree slug to "fix-auth"
    // Invoke test handler with DryRunner
    // Assert "PUBGOLF_DB_PORT=<expected>" in recorded env
}
```

### 5. Config Default Tests — `config_test.go`

```go
func TestCLIConfig_SetDefaults(t *testing.T) {
    c := CLIConfig{ProjectName: "pubgolf"}
    c.setDefaults()
    assert.Equal(t, "pubgolf-devctrl", c.CLIName)
    assert.Equal(t, "pubgolf-api-server", c.ServerBinName)
    assert.Equal(t, "dev", c.DopplerEnvName)
}
```

### 6. Project Root Resolution Tests — `root_test.go`

```go
func TestResolveProjectRoot_FromRoot(t *testing.T) {
    // In a temp dir with go.mod, resolveProjectRoot should return CWD
}

func TestResolveProjectRoot_FromSubdir(t *testing.T) {
    // In a temp dir with go.mod, from a subdirectory, should find parent
}

func TestResolveProjectRoot_NotFound(t *testing.T) {
    // In /tmp with no go.mod anywhere up, should return error
}
```

### 7. Error Classification Tests — `errors_test.go`

```go
func TestIsInfraError(t *testing.T) {
    tests := []struct{ err string; want bool }{
        {"address already in use", true},
        {"Cannot connect to the Docker daemon", true},
        {"connection refused", true},
        {"exit status 1", false},
        {"test failed", false},
    }
    for _, tt := range tests {
        assert.Equal(t, tt.want, isInfraError(fmt.Errorf(tt.err)))
    }
}
```

### 8. Error Simulation Tests — `runner_test.go`

```go
func TestDryRunner_ErrorFor(t *testing.T) {
    dr := &DryRunner{
        ErrorFor: map[string]error{
            "golangci-lint": fmt.Errorf("lint failed"),
        },
    }
    err := checkGo(context.Background(), dr)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "lint failed")
}
```

### 9. CLAUDE.md Updates

Add to the task runner table:

```markdown
| Inspect worktree status | `pubgolf-devctrl status` | auto-approved |
| Inspect all worktrees | `pubgolf-devctrl status --all` | auto-approved |
| Check dev environment | `pubgolf-devctrl doctor` | auto-approved |
| Preview cleanup | `pubgolf-devctrl clean` | auto-approved |
| Dry-run any command | `pubgolf-devctrl --dry-run <cmd>` | auto-approved |
| Force cleanup | `pubgolf-devctrl clean --force` | requires approval |
```

Add agent guidance section:

```markdown
## Exit Codes
- Exit 1: test/lint failure (debug your code).
- Exit 2: infrastructure failure (check services, not your code).

## Worktree Port Collisions
If you encounter "address already in use", set `PUBGOLF_PORT_OFFSET` to a
number between 1 and 499 and retry:
    export PUBGOLF_PORT_OFFSET=42
```

### 10. Verify

- `go test ./tools/...` — all tests pass
- `pubgolf-devctrl check:go` — no lint issues

## Acceptance Criteria

- [ ] Command construction tests for all command handlers
- [ ] Flag propagation tests for `--verbose`, `--coverage`, `--local`
- [ ] Worktree isolation tests (project name, port, slug injection)
- [ ] Config default tests
- [ ] Project root resolution tests
- [ ] Error classification tests
- [ ] Error simulation tests via DryRunner.ErrorFor
- [ ] CLAUDE.md updated with new commands and agent guidance
- [ ] All tests run with `go test ./tools/...` (no external deps)
- [ ] `pubgolf-devctrl check:go` passes
