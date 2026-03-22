# Plan 01: Runner Interface + Dry-Run

**Depends on:** nothing
**Branch:** `devctrl/01-runner-interface`
**PR scope:** `tools/lib/cmd/runner.go`, modifications to all command files to
thread Runner

## Objective

Introduce the `Runner` and `Process` interfaces, `ExecRunner` and `DryRunner`
implementations, and the `--dry-run` persistent flag. Refactor all command
files to route subprocess invocations through Runner.

## Steps

### 1. Create `tools/lib/cmd/runner.go`

Define the types:

```go
type Cmd struct {
    Name   string
    Args   []string
    Dir    string
    Env    []string
    Stdin  io.Reader
    Stdout io.Writer
    Stderr io.Writer
}

type Process interface {
    Wait() error
    Stop()
}

type Runner interface {
    Run(ctx context.Context, cmd Cmd) error
    Start(ctx context.Context, cmd Cmd) (Process, error)
}
```

Add a `String()` method on `Cmd` that shell-quotes arguments for copy-pasteable
output.

### 2. Implement `ExecRunner`

- `Run`: build `exec.Cmd`, merge env (explicit dedup: parse `os.Environ()` into
  map, overlay `Cmd.Env`, flatten), wire I/O, log before exec, call `cmd.Run()`.
- `Start`: same setup but call `cmd.Start()` with `Setpgid: true`. Return an
  `execProcess` that wraps `*exec.Cmd`. `Stop()` sends SIGINT to process group.
  `Wait()` calls `cmd.Wait()`.
- Log format: `[pubgolf-devctrl] exec: <cmd.String()>` with `dir:` and `env:`
  lines if non-empty.

### 3. Implement `DryRunner`

- `Recorded []Cmd` slice.
- `ErrorFor map[string]error` — keyed on `Cmd.Name`.
- `Run`: append to Recorded, log dry-run format, return ErrorFor[name] or nil.
- `Start`: append to Recorded, return a `dryProcess` (Wait returns nil, Stop
  is no-op).
- Log format matches the spec's sequence-numbered output. For this PR, omit
  the sequence numbers (they depend on `runPar` refactor in 05-errors) — just
  log `[pubgolf-devctrl] dry-run: <cmd.String()>`.

### 4. Add `--dry-run` flag to `rootCmd`

In `root.go`:
- Add `--dry-run` persistent flag on `rootCmd`.
- In `Execute()`, after parsing flags, set the package-level `runner` to
  `ExecRunner{}` or `DryRunner{}` based on the flag.
- Skip `checkVersion()` in dry-run mode.

### 5. Thread Runner through all command files

For each file, extract the `exec.CommandContext` calls into the extracted
functions and add a `runner Runner` parameter:

- **`check.go`**: `checkGo(ctx, runner)`, `checkWeb(ctx, runner)`,
  `checkProto(ctx, runner)`. Each builds a `Cmd` and calls `runner.Run()`.
- **`generate.go`**: `generateProto(ctx, runner)`, `generateSQLc(ctx, runner)`,
  `generateEnum(ctx, runner)`, `generateMock(ctx, runner)`,
  `generateInterface(ctx, runner)`.
- **`test.go`**: The test command builds a `Cmd` for `doppler run ... go test`.
  Use `runner.Run()`. For e2e, use `runner.Start()` returning a `Process`.
- **`run.go`**: `dopplerGoRun` → use `runner.Start()`. `dopplerDockerRun` →
  use `runner.Run()`. Stop function now calls `process.Stop()`.
- **`stop.go`**: `dopplerDockerStop` → use `runner.Run()`.
- **`install.go`**: `installWithGolang`, `installWithHomebrew` → `runner.Run()`.
- **`update.go`**: `go install` call → `runner.Run()`.
- **`migrate.go`**: `migrateCreate`'s `migrate create` call → `runner.Run()`.

The cobra `Run` closures call the extracted functions with the package-level
`runner`.

**Important:** `readDopplerVars` in `config.go` also does `exec.CommandContext`.
Route this through Runner too. In dry-run mode it returns an empty map (causing
defaults to be used).

### 6. Verify

- `pubgolf-devctrl check:go` (should still work identically)
- `go test ./tools/...` (new runner_test.go with at least one smoke test)
- `pubgolf-devctrl --dry-run check go` (should print the golangci-lint command
  without executing)

## Acceptance Criteria

- [ ] `Runner` and `Process` interfaces defined and exported
- [ ] `ExecRunner` passes existing manual smoke tests (commands still work)
- [ ] `DryRunner` records all commands
- [ ] `--dry-run` flag works on all subcommands
- [ ] `readDopplerVars` routed through Runner
- [ ] At least one test in `runner_test.go` verifying DryRunner recording
- [ ] `pubgolf-devctrl check:go` passes
