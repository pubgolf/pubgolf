# pubgolf-devctrl: Runner Interface, Concurrent Worktree Support, and Robustness

**Date:** 2026-03-16
**Status:** Draft (v2 — incorporates Go, DevTools, SRE, and Agent Workflow reviews)
**Related PR:** https://github.com/pubgolf/pubgolf/pull/575 (worktree-agnostic dbtest + DB stability)

## Context

`pubgolf-devctrl` is the project's task runner wrapping golangci-lint, buf, sqlc,
mockery, doppler, docker-compose, and Go test. It has zero tests. All subprocess
invocations go through bare `exec.CommandContext` calls scattered across cobra
command handlers, making the tool opaque to debug and impossible to test without
side effects. Additionally, Doppler is tightly coupled into every command path —
you cannot run tests, start servers, or even stop Docker containers without it.

The project is moving toward a workflow where **multiple agents work concurrently
in separate git worktrees**. This surfaces four categories of problems:

1. **Linting on main bleeds into worktrees.** `golangci-lint run ./api/... ./tools/...`
   recursively includes `.worktrees/` contents, causing false lint failures on
   main when a worktree has in-flight changes.
2. **No way to validate devctrl behavior without executing real commands.** An agent
   debugging a devctrl issue must run the real tool and observe side effects,
   rather than inspecting what commands would be dispatched.
3. **Concurrent worktree runs collide on shared resources.** Docker container
   names, ports, database namespaces, and volume mounts are all hardcoded,
   making it impossible to run `devctrl run` or `devctrl test` from multiple
   worktrees simultaneously.
4. **Doppler coupling prevents autonomous operation.** Agents in worktrees may
   not have Doppler configured, and every command path currently requires it.

PR #575 separately addresses the `dbtest.MigrationDir()` hardcoded path issue
and DB connection pool stability. This spec builds on that work.

## Goals

1. Introduce a `Runner` interface that all subprocess invocations route through.
2. Add a `--dry-run` flag with structured output for agent debugging.
3. Make linting and code-generation tools worktree-safe on main.
4. Make Docker services, test databases, and API servers safe for concurrent
   worktree execution.
5. Decouple Doppler from devctrl's core command paths.
6. Add agent-oriented commands: `status`, `clean`, `doctor`.
7. Replace `guard`/`log.Fatalf` with classified error handling and distinct
   exit codes.
8. Fix confirmed brittleness issues found during audit.
9. Add tests exercising the dry-run path to validate command construction.

## Non-Goals (Explicitly Out of Scope)

These are items that were considered and deliberately excluded. Each has a
rationale explaining why.

1. **Rewriting the cobra command tree structure.** The existing command
   hierarchy (`check`, `generate`, `test`, `run`, etc.) is fine. The refactor
   changes internals (Runner threading, error handling) not the CLI surface.

2. **Multi-machine or CI-level parallelism.** This spec targets local
   developer/agent tooling on a single machine. CI runs are a separate concern
   with different isolation requirements (containers, ephemeral VMs).

3. **Automatic cleanup on `git worktree remove`.** Git worktree lifecycle
   hooks are fragile and platform-dependent. Explicit `devctrl clean` is more
   reliable and debuggable.

4. **Agent-to-agent service discovery / registry.** A shared registry file
   (`data/worktree-services.json`) for cross-worktree port discovery was
   considered. Deferred because `devctrl status --all` provides the same
   information on-demand, and cross-agent API calls are not a current workflow.

5. **`fsnotify` migration.** Replacing `radovskyb/watcher` (polling) with
   `fsnotify` (OS-level events) would reduce CPU overhead with many worktrees.
   Deferred because the watcher is only active during `--watch` mode and the
   current 100ms polling is acceptable for the near term.

6. **Shared Postgres `max_connections` tuning.** Concurrent worktrees in
   shared-postgres mode could exhaust the default `max_connections=100`.
   Deferred because embedded mode (the default) avoids this entirely, and
   shared mode is an opt-in power-user path.

7. **`devctrl logs` command.** Tailing Docker container logs for the current
   worktree's project (wrapping `docker logs <project>-api-db-1`). Useful but
   low priority — agents can use `docker logs` directly if needed, and
   `devctrl status` shows the container names.

8. **Web dev server port offsetting.** The SvelteKit dev server (Vite port
   5173) and Playwright preview (port 4173) are started manually via
   `npm run dev`, not through devctrl. Offsetting these would require changes
   to the web-app's vite config and devctrl's `run` command. Documented as a
   known limitation — concurrent web dev servers require manual `--port` flags.

## Implementation Strategy

This spec is designed for **parallel agent execution via independent PRs**.
The work decomposes into modules with minimal cross-dependencies:

- **Runner interface + dry-run** — standalone, no dependency on other changes.
- **Env provider decoupling** — depends on Runner (uses `Cmd.Env`).
- **Worktree identity + port offset** — standalone.
- **Docker isolation** — depends on worktree identity + env provider.
- **Error handling refactor** — depends on Runner (functions return errors).
- **New commands** (`status`, `clean`, `doctor`) — depend on worktree identity.
- **Robustness fixes** — each is independent and can be its own PR.
- **Test suite** — depends on Runner + DryRunner being available.
- **CLAUDE.md updates** — depends on new commands being merged.

The implementation plan (next step) will identify which of these can be
parallelized across agents and which must be sequenced.

---

## Design

### 1. Runner Interface

New file: `tools/lib/cmd/runner.go`

```go
// Cmd describes a subprocess invocation.
type Cmd struct {
    Name   string
    Args   []string
    Dir    string    // working directory; empty uses projectRoot
    Env    []string  // additional env vars; nil inherits parent env
    Stdin  io.Reader // nil = os.Stdin
    Stdout io.Writer // nil = os.Stdout
    Stderr io.Writer // nil = os.Stderr
}

// Process represents a running subprocess started via Runner.Start.
type Process interface {
    // Wait blocks until the process exits and returns its error.
    Wait() error
    // Stop sends SIGINT to the process group for graceful shutdown.
    Stop()
}

// Runner executes or records subprocess invocations.
type Runner interface {
    // Run executes (or records) a single command and waits for completion.
    Run(ctx context.Context, cmd Cmd) error
    // Start begins a long-running process and returns a Process handle.
    Start(ctx context.Context, cmd Cmd) (Process, error)
}
```

The `Process` interface replaces the previous `stop func()` return, providing
both `Wait()` (to detect natural exit and trigger shutdown) and `Stop()` (to
request graceful shutdown). This matches the existing patterns in `run.go` and
`test.go` where both capabilities are needed.

The `Env` field uses explicit deduplication when merging with `os.Environ()`:
parse parent env into a map, overlay `Cmd.Env`, flatten back to a slice. This
avoids relying on platform-specific last-wins behavior of `execve`.

**`ExecRunner`** — production implementation:
- Calls `exec.CommandContext` with the fields from `Cmd`.
- Logs before execution (see section 3 for format).
- Wires stdout/stderr/stdin from `Cmd`, defaulting to `os.Stdout`/`os.Stderr`/`os.Stdin`.
- Merges `Cmd.Env` with `os.Environ()` using explicit deduplication.
- `Start()` calls `cmd.Start()` with `Setpgid: true` and returns an
  `execProcess` implementing `Process`. `Stop()` sends SIGINT (not SIGKILL —
  this is an intentional behavioral fix; the existing code at `run.go:156`
  sends SIGKILL despite its comment saying SIGINT).

**`DryRunner`** — test and dry-run implementation:
- Appends each `Cmd` to a `[]Cmd` slice (`DryRunner.Recorded`).
- Logs in the format specified by section 3.
- Returns `nil` (success) by default. Tests can set `DryRunner.ErrorFor` map
  to simulate failures for specific command names.
- `Start()` records the command and returns a no-op `Process` (Wait returns
  nil immediately, Stop is a no-op).

### 2. Threading the Runner

A package-level `runner Runner` variable is set in `Execute()` based on the
`--dry-run` persistent flag on `rootCmd`. All command handlers receive `runner`
implicitly through the package scope (same pattern as `config` today).

Each extracted function (`checkGo`, `checkWeb`, `checkProto`, `generateProto`,
`generateSQLc`, etc.) gains a `runner Runner` parameter. The cobra `Run`
closures pass the package-level `runner`.

**`readDopplerVars`** (in `config.go`) is also routed through `Runner`. In
dry-run mode, it records the `doppler secrets` invocation and returns an empty
map, causing downstream config to use defaults. This ensures `--dry-run` has
no side effects.

Migration commands (`migrateUp`, `migrateDown`) that use the `migrate` library
directly (not subprocess) are unaffected by the Runner — only `migrateCreate`'s
shell-out to `migrate create` is routed through `Runner`.

### 3. Dry-Run Output

```
pubgolf-devctrl --dry-run check go
```

Output format includes sequence numbers, parallel/sequential annotation, and
injected environment variables:

```
[pubgolf-devctrl] dry-run (1/1): golangci-lint run --exclude-dirs '\.worktrees' ./api/... ./tools/...
  dir: /Users/dev/Sites/pubgolf
```

For multi-step commands like `generate`:

```
[pubgolf-devctrl] dry-run (1/5, parallel): buf generate --template buf.gen.dev.yaml
[pubgolf-devctrl] dry-run (2/5, parallel): sqlc generate --file api/internal/db/sqlc.yaml
[pubgolf-devctrl] dry-run (3/5, parallel): enumer -sql -transform snake-upper -type ScoringCategory ./api/internal/lib/models
[pubgolf-devctrl] dry-run (4/5, sequential): ifacemaker --struct Querier ...
  dir: /Users/dev/Sites/pubgolf
[pubgolf-devctrl] dry-run (5/5, sequential): mockery --dir api/internal/lib/dao/ ...
```

For commands with env var injection (worktree context):

```
[pubgolf-devctrl] dry-run (1/1): go test ./api/...
  dir: /Users/dev/Sites/pubgolf/.worktrees/fix-auth
  env: PUBGOLF_DB_PORT=5437 PUBGOLF_WORKTREE_SLUG=fix-auth PUBGOLF_PORT=5005
```

The `Cmd` struct has a `String()` method that shell-quotes arguments correctly
so output is copy-pasteable. The `--dry-run` flag is persistent on `rootCmd`,
available to all subcommands. In dry-run mode, `checkVersion()` is skipped.

### 4. Decoupling Doppler

Currently Doppler wraps every command that needs secrets or environment config.
This creates a hard dependency: agents without Doppler configured cannot run
any devctrl command.

#### 4a. Environment Provider Interface

New file: `tools/lib/cmd/envprovider.go`

```go
// EnvProvider resolves environment variables for subprocess execution.
type EnvProvider interface {
    // Env returns the environment variables for the given context.
    // The project and config parameters correspond to Doppler's concepts
    // but are provider-agnostic.
    Env(ctx context.Context, project, config string) ([]string, error)
}
```

**`DopplerProvider`** — wraps commands in `doppler run --project P --config C --`:
- This is the current behavior, extracted into an interface implementation.
- For test commands, the provider wraps the entire `go test` invocation.
- For Docker commands, the provider wraps `docker-compose`.

**`EnvFileProvider`** — reads from `.env` files or the process environment:
- Looks for `.env.{config}` (e.g., `.env.dev`, `.env.test`) in the project root.
- Falls back to the current process environment.
- This allows agents to work without Doppler by providing env vars directly.

**`AutoProvider`** — default, tries Doppler first, falls back to env files:
- Checks if `doppler` is in PATH and configured (`doppler status`).
- If available, uses `DopplerProvider`.
- If not, uses `EnvFileProvider` with a log message:
  `"[pubgolf-devctrl] Doppler not available, using environment variables."`

The provider is selected at startup in `Execute()` and stored as a package-level
variable alongside `runner` and `config`. A `--env-provider` flag allows explicit
selection: `doppler`, `env`, or `auto` (default).

#### 4b. Removing Doppler from Command Paths

Commands that currently wrap everything in `doppler run ... --` are refactored:

- **`dopplerDockerRun` → `dockerRun`**: Uses `envProvider.Env()` to get the
  environment, then passes it via `Cmd.Env` to `docker-compose` directly.
- **`dopplerDockerStop` → `dockerStop`**: Same — no secrets needed for `down`.
- **`dopplerGoRun` → `goRun`**: Gets env from provider, passes to `go run`.
- **Test command**: Gets env from provider, passes to `go test`.

The `readDopplerVars` / `getDatabaseURL` functions in `config.go` are refactored
to use `EnvProvider` instead of shelling out to `doppler secrets` directly.

### 5. Concurrent Worktree Isolation

Multiple worktrees must be able to run devctrl commands simultaneously without
resource collisions.

#### 5a. Worktree Identity

At startup in `Execute()`, after resolving the project root, compute a
**worktree slug**: a short, normalized identifier derived from the worktree
directory name.

```go
func worktreeSlug() (string, error) {
    // Check if we're in a worktree (not the main working tree)
    topLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
    if err != nil {
        return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
    }
    topLevelAbs, _ := filepath.Abs(strings.TrimSpace(string(topLevel)))

    commonDir, err := exec.Command("git", "rev-parse", "--git-common-dir").Output()
    if err != nil {
        return "", fmt.Errorf("git rev-parse --git-common-dir: %w", err)
    }
    commonDirAbs, _ := filepath.Abs(strings.TrimSpace(string(commonDir)))

    // If commonDir is outside topLevel, we're in a worktree.
    if !strings.HasPrefix(commonDirAbs, topLevelAbs) {
        raw := filepath.Base(topLevelAbs)
        return normalizeSlug(raw), nil
    }

    return "", nil // main working tree — use defaults
}
```

Key differences from v1:
- Returns `(string, error)` — git failures are errors, not silent fallbacks.
  An error causes devctrl to exit with a clear message rather than silently
  using the main tree's resources.
- Both paths resolved to absolute via `filepath.Abs` before comparison (fixes
  `--git-common-dir` returning relative paths in some git versions).
- Slug is normalized via `normalizeSlug()`.

**Slug normalization:** Lowercase, replace non-alphanumeric characters with
hyphens, collapse consecutive hyphens, truncate to 20 characters. If truncation
occurs, append a 6-character FNV hash suffix for uniqueness:

```
fix-auth                    → fix-auth
feature/minio-s3-support    → minio-s3-support
issue-1234-some-very-long-description → issue-1234-some-very-a3f2b1
```

This keeps Docker container names (64-char limit) and PostgreSQL identifiers
(63-byte limit) well within bounds.

#### 5b. Port Offset

```go
func worktreePortOffset() int {
    slug, _ := worktreeSlug()
    if slug == "" {
        return 0
    }
    // Check for explicit override first
    if v := os.Getenv("PUBGOLF_PORT_OFFSET"); v != "" {
        offset, err := strconv.Atoi(v)
        if err == nil && offset > 0 && offset < 500 {
            return offset
        }
    }
    // Stable hash of slug
    h := fnv.New32a()
    h.Write([]byte(slug))
    return int(h.Sum32()%500) + 1
}
```

Changes from v1:
- **`PUBGOLF_PORT_OFFSET` env var override** — deterministic escape hatch for
  hash collisions. An agent can `export PUBGOLF_PORT_OFFSET=42` and retry
  without recreating the worktree.
- **Range widened from mod 99 to mod 500** — collision probability drops from
  ~10% at 5 worktrees to well under 1%. Base ports (5432, 5000) plus offsets
  up to 500 stay well below contested port ranges.
- A `--port-offset <int>` flag on `devctrl run` and `devctrl test` provides
  the same override for interactive use. The flag sets `PUBGOLF_PORT_OFFSET`
  in the process environment so it propagates to all subcommands.

The offset applies to:
- DB port: `5432 + offset`
- API port: `5000 + offset`

**Web dev server ports** (Vite default 5173, Playwright preview 4173) are not
offset by devctrl — these are run manually via `npm run dev`, not through
devctrl. This is a **known limitation**: concurrent web dev servers in different
worktrees require manual port configuration via `--port` flags to Vite.

#### 5c. Docker Service Isolation

Pass a per-worktree project name to docker-compose in both the `up` and `down`
paths (`dockerRun` and `dockerStop`):

```go
projectName := "pubgolf"
if slug, _ := worktreeSlug(); slug != "" {
    projectName = "pubgolf-" + slug
}
```

Docker Compose's `--project-name` flag namespaces all container names, networks,
and volumes automatically.

**Volume isolation:** Devctrl sets `PUBGOLF_DB_HOST_DATA_PATH` to a
worktree-specific path: `./data/postgres-{slug}` for worktrees, `./data/postgres`
for main.

**Coverage directory:** `./data/go-test-coverage-{slug}` for worktrees to
prevent concurrent `--coverage` runs from racing on the same directory.

#### 5d. Test Database Isolation (Shared Mode)

Modify `dbtest.NewConn()` to include a worktree discriminator in the database
name when `PUBGOLF_WORKTREE_SLUG` is set. Devctrl's test command injects this
env var.

**`PUBGOLF_SHARED_DB_URL` must reflect the offset port.** When devctrl runs
tests in shared-postgres mode from a worktree, it constructs the shared DB URL
using the offset port rather than reading a static URL from Doppler. The env
provider computes this: base connection details from Doppler/env, port from
the offset calculation.

**Embedded mode** (default, no `-shared-postgres`): Already safe. Uses
`freeport.GetFreePort()` for dynamic port allocation and isolated temp
directories per test package.

#### 5e. Pre-Flight Health Check

Before running tests, devctrl performs a quick infrastructure check:

1. If shared-postgres mode: TCP dial the expected DB port.
2. If Docker services are expected: check container health via
   `docker inspect --format '{{.State.Health.Status}}'`.

If either check fails, exit with code 2 (infrastructure failure) and a message:

```
[pubgolf-devctrl] ERROR: Infrastructure check failed
  DB container is not running on port 5437.
  Run 'pubgolf-devctrl run bg' to start it.
  This is an infrastructure issue, not a code problem.
```

Exit code 2 is distinct from exit code 1 (test/lint failure). Agents can be
taught (in CLAUDE.md) that exit code 2 means "don't debug your code."

### 6. Worktree Exclusions (Linting/Codegen on Main)

**golangci-lint** (`check.go`):
Add `--exclude-dirs` flag with pattern `\.worktrees`:
```go
"golangci-lint", "run", "--exclude-dirs", `\.worktrees`, "./api/...", "./tools/..."
```

**buf** (`check.go`, `generate.go`):
Add `.worktrees` to `excludes` in `buf.yaml`.

**File watchers** (`root.go`):
Centralized ignore list applied by `watch()`:
```go
var ignoredDirs = []string{".worktrees", "node_modules", "vendor", "data", ".git"}
```
This is defensive and avoids per-watcher ignore rules as new watchers are added.

### 7. Error Handling

#### 7a. Replace `guard` with Classified Errors

The current `guard` function calls `log.Fatalf`, which:
- Always exits with code 1 (no distinction between code and infrastructure errors).
- Kills the process mid-goroutine in `runPar`, preventing cleanup.
- Produces messages like `"execute docker-compose up ... command: exit status 1"` —
  useless for agent self-correction.

Replace `guard` with an error-return pattern. All extracted functions return
`error`. The cobra `Run` closures handle errors with classified exit:

```go
func classifyAndExit(err error) {
    if err == nil {
        return
    }
    if isInfraError(err) {
        log.Printf("[pubgolf-devctrl] ERROR: Infrastructure failure\n  %s", err)
        os.Exit(2)
    }
    log.Printf("[pubgolf-devctrl] ERROR: %s", err)
    os.Exit(1)
}
```

`isInfraError` checks for known patterns in stderr: `"address already in use"`,
`"Cannot connect to the Docker daemon"`, `"permission denied"`,
`"connection refused"`.

#### 7b. `runPar` Error Collection

`runPar` currently calls `guard` in goroutines (which calls `os.Exit`). Refactor
to collect errors from all parallel commands and report them together:

```go
func runPar(ctx context.Context, runner Runner, fns ...func(context.Context, Runner) error) error {
    var mu sync.Mutex
    var errs []error
    var wg sync.WaitGroup
    for _, fn := range fns {
        wg.Add(1)
        go func(f func(context.Context, Runner) error) {
            defer wg.Done()
            if err := f(ctx, runner); err != nil {
                mu.Lock()
                errs = append(errs, err)
                mu.Unlock()
            }
        }(fn)
    }
    wg.Wait()
    return errors.Join(errs...)
}
```

#### 7c. Port Collision Error Messages

When a port bind failure is detected:

```
[pubgolf-devctrl] ERROR: port 5437 is already in use (PUBGOLF_DB_PORT)
  Worktree "fix-auth" was assigned port offset 5 (base 5432 + 5 = 5437).
  Likely cause: another worktree has the same port offset, or a previous
  instance was not stopped cleanly.
  To resolve:
    1. Run: pubgolf-devctrl stop          (in the conflicting worktree)
    2. Or:  export PUBGOLF_PORT_OFFSET=42 && pubgolf-devctrl run bg
    3. Or:  pubgolf-devctrl clean --force  (from main, removes orphaned resources)
```

### 8. Robustness Fixes

#### 8a. `stopCmd` Description at Init Time

Hardcode the description string. There is only one project using this tool.

#### 8b. `migrateCreate` Stderr Parsing

Guard against empty output. Filter lines to those containing the migration
directory prefix. Log a warning if no files were found.

#### 8c. Project Root Resolution

In `Execute()`, resolve the project root by walking up to find `go.mod`. Store
as `projectRoot string` (package-level). All relative paths in command
construction use `filepath.Join(projectRoot, ...)` — no `os.Chdir`.

Log clearly: `"[pubgolf-devctrl] Resolved project root: /path/to/root"`.
If not found: exit with `"pubgolf-devctrl must be run from within the project directory"`.

#### 8d. `beginShutdown` / Shutdown Safety

Replace with `sync.Once`-guarded `triggerShutdown()`:

```go
var shutdownOnce sync.Once
func triggerShutdown() {
    shutdownOnce.Do(func() { close(shuttingDown) })
}
```

Signal handler in `PersistentPreRun` calls `triggerShutdown()` directly,
eliminating the `beginShutdown` intermediary channel.

**Shutdown order:** API server stopped first (SIGINT via `Process.Stop()`),
then Docker containers remain running (they are `--detach`ed background
services). `devctrl stop` is required for full teardown.

#### 8e. `checkVersion` Absolute Path

Use `filepath.Join(projectRoot, "tools")` for the `dirhash.HashDir` call.

#### 8f. SIGKILL → SIGINT

The existing `dopplerGoRun` stop function sends `SIGKILL` despite its comment
saying `SIGINT`. The `ExecRunner.Start()` implementation sends `SIGINT` for
graceful shutdown. This is an intentional behavioral fix.

### 9. New Commands

#### 9a. `devctrl status`

Read-only introspection. Safe for autonomous agent use.

```
$ pubgolf-devctrl status
Worktree:        fix-auth
Slug:            fix-auth
Port offset:     5
Docker project:  pubgolf-fix-auth
DB port:         5437
API port:        5005
DB volume:       ./data/postgres-fix-auth
Docker status:   running (api-db: healthy)
Env provider:    doppler (project: pubgolf-api-server, config: dev)
```

```
$ pubgolf-devctrl status --all
WORKTREE         SLUG          OFFSET  DB_PORT  API_PORT  DOCKER
(main)           -             0       5432     5000      running
fix-auth         fix-auth      5       5437     5005      running
add-leaderboard  add-leader..  23      5455     5023      stopped
```

`--all` enumerates worktrees via `git worktree list`, computes each slug and
offset, and checks Docker status for each project name.

#### 9b. `devctrl clean`

Discovers and removes orphaned worktree resources. Runs from any worktree or
the main working tree.

**Defaults to dry-run.** Must pass `--force` to actually delete.

```
$ pubgolf-devctrl clean
[pubgolf-devctrl] Dry run — pass --force to remove these resources:
  [orphan] docker project: pubgolf-old-branch (containers: 2, stopped)
  [orphan] data directory: data/postgres-old-branch (1.2 GB)

$ pubgolf-devctrl clean --force
[pubgolf-devctrl] Removed docker project: pubgolf-old-branch
[pubgolf-devctrl] Removed data directory: data/postgres-old-branch
```

**Two-pass algorithm:**
1. **Docker pass** (skipped if daemon unavailable, with warning): list all
   Docker Compose projects matching `pubgolf-*`, cross-reference against
   `git worktree list`. Orphaned projects get `docker-compose down --volumes`.
2. **Filesystem pass** (always runs): scan `data/postgres-*` and
   `data/go-test-coverage-*` directories, cross-reference against active
   worktree slugs. Orphaned directories are removed.

When Docker is unavailable:
```
[pubgolf-devctrl] WARNING: Docker daemon is not available. Skipping container cleanup.
  Start Docker and re-run to clean up Docker resources.
```

#### 9c. `devctrl stop` Enhancements

- From a worktree: stops that worktree's Docker project.
- From main: stops the main project.
- `devctrl stop --all`: stops services across all worktree projects.
- `devctrl stop --remove-data`: also removes this worktree's data directories
  (Postgres volume, coverage). For explicit teardown when done with a worktree.

**Post-start reminder:** When `devctrl run bg` starts services from a worktree:
```
[pubgolf-devctrl] Started services for worktree 'fix-auth' (DB: port 5437).
  Run 'pubgolf-devctrl stop' before removing this worktree.
```

#### 9d. `devctrl doctor`

Comprehensive environment check. Safe for autonomous use.

```
$ pubgolf-devctrl doctor
[pubgolf-devctrl] Environment check:
  Go:             1.26.0    ✓
  Docker:         running   ✓
  Doppler:        configured (project: pubgolf-api-server)  ✓
  golangci-lint:  2.11.3    ✓
  buf:            installed  ✓
  sqlc:           1.24.0    ✓
  mockery:        2.42.2    ✓
  Node.js:        22.x      ✓
  npm:            installed  ✓
```

Reports all issues at once rather than letting agents discover them one
command at a time.

### 10. Test Strategy

New file: `tools/lib/cmd/runner_test.go` (and additional `*_test.go` files
as needed for each command file).

Tests construct a `DryRunner`, set up a `CLIConfig`, and invoke command
functions directly, then assert on `DryRunner.Recorded`:

```go
func TestCheckGo_CommandConstruction(t *testing.T) {
    dr := &DryRunner{}
    err := checkGo(context.Background(), dr)
    require.NoError(t, err)
    require.Len(t, dr.Recorded, 1)

    cmd := dr.Recorded[0]
    assert.Equal(t, "golangci-lint", cmd.Name)
    assert.Contains(t, cmd.Args, "--exclude-dirs")
}
```

Test categories:

- **Command construction tests** — verify each command handler produces the
  expected subprocess invocation (name, args, working dir, env vars).
- **Flag propagation tests** — verify `--verbose`, `--coverage`, `--local`,
  `--watch` flags modify the constructed commands correctly.
- **Worktree isolation tests** — verify that when a worktree slug is set,
  Docker project names are suffixed, ports are offset, and
  `PUBGOLF_WORKTREE_SLUG` is injected into test environments.
- **Config default tests** — verify `CLIConfig.setDefaults()` behavior.
- **Project root resolution tests** — verify detection from subdirectories.
- **Port offset tests** — verify determinism, env var override, and that
  different slug values produce non-colliding offsets.
- **Slug normalization tests** — verify truncation, hash suffix, and
  character replacement.
- **Error classification tests** — verify `isInfraError` correctly classifies
  known error patterns.
- **Env provider tests** — verify `AutoProvider` fallback behavior,
  `EnvFileProvider` parsing.
- **Error simulation tests** — set `DryRunner.ErrorFor` to simulate command
  failures and verify error handling paths.

Tests do NOT require Doppler, a database, or any external tools. They run with
plain `go test ./tools/...`.

### 11. Files Changed

| File | Change |
|------|--------|
| `tools/lib/cmd/runner.go` | New — `Runner`, `Process` interfaces, `ExecRunner`, `DryRunner` |
| `tools/lib/cmd/runner_test.go` | New — dry-run, worktree isolation, error classification tests |
| `tools/lib/cmd/envprovider.go` | New — `EnvProvider` interface, `DopplerProvider`, `EnvFileProvider`, `AutoProvider` |
| `tools/lib/cmd/envprovider_test.go` | New — provider fallback and parsing tests |
| `tools/lib/cmd/worktree.go` | New — `worktreeSlug()`, `normalizeSlug()`, `worktreePortOffset()` |
| `tools/lib/cmd/worktree_test.go` | New — slug normalization, port offset, identity tests |
| `tools/lib/cmd/status.go` | New — `devctrl status` and `status --all` |
| `tools/lib/cmd/clean.go` | New — `devctrl clean` with dry-run default and `--force` |
| `tools/lib/cmd/doctor.go` | New — `devctrl doctor` environment check |
| `tools/lib/cmd/errors.go` | New — `classifyAndExit`, `isInfraError`, exit code constants |
| `tools/lib/cmd/root.go` | `--dry-run` flag, `--env-provider` flag, project root resolution, `triggerShutdown()`, thread runner |
| `tools/lib/cmd/check.go` | Thread runner, add `--exclude-dirs` to golangci-lint |
| `tools/lib/cmd/generate.go` | Thread runner through all generate functions |
| `tools/lib/cmd/test.go` | Thread runner, `triggerShutdown()`, inject worktree env, pre-flight health check, `--port-offset` flag |
| `tools/lib/cmd/run.go` | Thread runner, `Process` interface, `triggerShutdown()`, worktree port/project isolation, SIGINT fix, post-start reminder, `--port-offset` flag |
| `tools/lib/cmd/stop.go` | Thread runner, hardcode description, worktree project name, `--all`, `--remove-data` |
| `tools/lib/cmd/install.go` | Thread runner |
| `tools/lib/cmd/update.go` | Thread runner, absolute path for dirhash |
| `tools/lib/cmd/migrate.go` | Thread runner for `migrateCreate`, defensive stderr parsing |
| `tools/lib/cmd/config.go` | `projectRoot` resolution, refactor `getDatabaseURL` to use `EnvProvider` |
| `api/internal/lib/dbtest/dbtest.go` | Read `PUBGOLF_WORKTREE_SLUG` for shared-mode DB name isolation |
| `buf.yaml` | Add `.worktrees` to excludes |
| `CLAUDE.md` | Add `status`, `doctor`, `clean --dry-run` to pre-approved; add `clean --force` to requires-approval; document exit codes and `PUBGOLF_PORT_OFFSET` |

### 12. CLAUDE.md Updates

New pre-approved commands:
- `pubgolf-devctrl status` (read-only)
- `pubgolf-devctrl status --all` (read-only)
- `pubgolf-devctrl doctor` (read-only)
- `pubgolf-devctrl clean` (dry-run by default, read-only)
- `pubgolf-devctrl --dry-run <anything>` (read-only)

Requires approval:
- `pubgolf-devctrl clean --force` (destroys Docker resources)

Agent guidance:
- Exit code 1 = test/lint failure (debug your code).
- Exit code 2 = infrastructure failure (don't debug your code; check services).
- If port collision: `export PUBGOLF_PORT_OFFSET=<number> && retry`.

### 13. Migration / Rollback

All changes are backward-compatible. The `--dry-run` flag is additive. The
`Runner` refactor changes internal structure but not external behavior. The
worktree isolation is a no-op when running from the main working tree (slug is
empty, offsets are zero). The `AutoProvider` defaults to Doppler when available,
preserving existing behavior.

Docker volumes created by worktrees (`data/postgres-{slug}`) persist after the
worktree is removed. `devctrl clean --force` handles cleanup.

If the refactor introduces regressions, reverting the commit restores the
previous behavior completely.

### 14. Known Limitations

1. **Web dev server ports** (Vite 5173, Playwright 4173) are not offset by
   devctrl. Concurrent web dev servers in different worktrees require manual
   `--port` flags to Vite.

2. **Embedded Postgres temp directories** (under `os.TempDir()`) are not
   cleaned by `devctrl clean`. They are cleaned by test cleanup functions but
   may leak if tests are killed. This is pre-existing behavior.

3. **`devctrl doctor` does not verify Doppler project/config access**, only
   that the CLI is installed. Verifying access would require a network call
   to Doppler's API.

4. **`radovskyb/watcher` polls at 100ms intervals.** With many concurrent
   worktrees, filesystem polling can become CPU-intensive. Migration to
   `fsnotify` (OS-level events) is a potential follow-up.
