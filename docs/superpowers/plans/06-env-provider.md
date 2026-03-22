# Plan 06: Env Provider Decoupling

**Depends on:** 01-runner-interface
**Branch:** `devctrl/06-env-provider`
**PR scope:** `tools/lib/cmd/envprovider.go`, `tools/lib/cmd/envprovider_test.go`,
modifications to `config.go`, `test.go`, `run.go`, `stop.go`

## Objective

Decouple Doppler from devctrl's core command paths by introducing a pluggable
`EnvProvider` interface. Commands that currently wrap everything in
`doppler run ... --` will instead get their environment from the provider and
pass it via `Cmd.Env`.

## Prerequisites

Plan 01 must be merged. The `Cmd.Env` field and `Runner.Run()` are needed for
passing environment variables to subprocesses.

## Steps

### 1. Create `tools/lib/cmd/envprovider.go`

```go
// EnvProvider resolves environment variables for subprocess execution.
type EnvProvider interface {
    // Env returns environment variables for a given project and config.
    Env(ctx context.Context, project, config string) ([]string, error)
}
```

#### `DopplerProvider`

- Calls `doppler secrets --project P --config C --json` via the `Runner`.
- Parses JSON output into `[]string` of `KEY=VALUE` pairs.
- This replaces the current `readDopplerVars` + `getDatabaseURL` pattern.
- For commands that currently wrap in `doppler run ... --`, the provider
  fetches the env and the command is invoked directly with `Cmd.Env`.

#### `EnvFileProvider`

- Looks for `.env.{config}` (e.g., `.env.dev`, `.env.test`) in `projectRoot`.
- Parses simple `KEY=VALUE` format (supports comments, blank lines).
- Falls back to the current process environment for missing keys.
- This allows agents to work without Doppler.

#### `AutoProvider`

- Default provider. Checks if `doppler` is available in PATH.
- If available: uses `DopplerProvider`.
- If not: uses `EnvFileProvider` with a log message:
  `"[pubgolf-devctrl] Doppler not available, using environment variables."`

### 2. Add `--env-provider` flag to `rootCmd`

In `root.go`:
- `--env-provider` persistent flag: `auto` (default), `doppler`, `env`.
- Set the package-level `envProvider` in `Execute()`.

### 3. Refactor command paths

#### `test.go`

Before:
```go
args := []string{"run", "--project", config.ServerBinName, "--config", "test", "--", "go", "test", ...}
tester := exec.CommandContext(ctx, "doppler", args...)
```

After:
```go
env, err := envProvider.Env(ctx, config.ServerBinName, "test")
if err != nil { return err }
cmd := Cmd{Name: "go", Args: []string{"test", "./api/..."}, Env: env}
return runner.Run(ctx, cmd)
```

#### `run.go`

`dopplerDockerRun` → `dockerRun`:
```go
env, err := envProvider.Env(ctx, project, envCfg)
cmd := Cmd{Name: "docker-compose", Args: [...], Env: env}
```

`dopplerDockerStop` → `dockerStop`:
```go
// docker-compose down doesn't need secrets; use minimal env or none
cmd := Cmd{Name: "docker-compose", Args: [...]}
```

`dopplerGoRun` → `goRun`:
```go
env, err := envProvider.Env(ctx, project, envCfg)
cmd := Cmd{Name: "go", Args: []string{"run", bin}, Env: env}
process, err := runner.Start(ctx, cmd)
```

#### `config.go`

`readDopplerVars` and `getDatabaseURL` are refactored to use `EnvProvider`:
```go
func getDatabaseURL(ctx context.Context, driver DBDriver, provider EnvProvider, project, env, prefix string, isMigrator bool) (string, error) {
    envVars, err := provider.Env(ctx, project, env)
    // Extract PUBGOLF_DB_* vars from envVars slice
    // Build URL from extracted vars with defaults
}
```

### 4. Create `tools/lib/cmd/envprovider_test.go`

- **EnvFileProvider parsing tests:** valid `.env` files, comments, blank lines,
  quoted values, missing files (falls back to os env).
- **AutoProvider fallback test:** mock `doppler` absence, verify fallback.
- **DopplerProvider JSON parsing test:** mock JSON output, verify env slice.

### 5. Verify

- `pubgolf-devctrl check:go` — passes
- `go test ./tools/...` — env provider tests pass
- `pubgolf-devctrl --env-provider=env test` — runs tests using process env
  (will fail if DB vars aren't set, but should not crash with a Doppler error)
- `pubgolf-devctrl --env-provider=doppler test` — existing behavior

## Acceptance Criteria

- [ ] `EnvProvider` interface defined
- [ ] `DopplerProvider`, `EnvFileProvider`, `AutoProvider` implemented
- [ ] All command paths use `envProvider.Env()` instead of `doppler run ... --`
- [ ] `--env-provider` flag on rootCmd
- [ ] `readDopplerVars`/`getDatabaseURL` refactored to use provider
- [ ] `docker-compose down` (stop) works without Doppler
- [ ] Tests verify provider parsing and fallback
- [ ] `pubgolf-devctrl check:go` passes
