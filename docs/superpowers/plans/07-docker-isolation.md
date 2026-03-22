# Plan 07: Docker + Test Database Isolation

**Depends on:** 01-runner-interface, 02-worktree-identity, 06-env-provider
**Branch:** `devctrl/07-docker-isolation`
**PR scope:** `tools/lib/cmd/run.go`, `tools/lib/cmd/stop.go`,
`tools/lib/cmd/test.go`, `api/internal/lib/dbtest/dbtest.go`

## Objective

Make Docker services, test databases, and API server ports safe for concurrent
worktree execution. This is the integration point where Runner, worktree
identity, and env provider come together.

## Prerequisites

Plans 01 (Runner), 02 (worktree identity), and 06 (env provider) must be
merged. This plan uses:
- `Runner` for subprocess execution
- `worktreeSlug()`, `worktreePortOffset()`, `worktreeDockerProject()`,
  `worktreeDataDir()` for isolation
- `EnvProvider` for secrets/environment resolution

## Steps

### 1. Docker Project Name — `run.go`, `stop.go`

In `dockerRun` and `dockerStop`, use `worktreeDockerProject()` for the
`--project-name` flag:

```go
projectName := worktreeDockerProject()
args := []string{
    "--file", filepath.Join(projectRoot, "infra/docker-compose.dev.yaml"),
    "--project-name", projectName,
    "up", "--detach",
}
```

**Critical:** Both `dockerRun` AND `dockerStop` must use the same project name.
The stop path is a separate code path — verify it uses `worktreeDockerProject()`.

### 2. Port Offset Injection — `run.go`, `test.go`

Compute the offset and inject port env vars:

```go
offset := worktreePortOffset()
env = append(env,
    fmt.Sprintf("PUBGOLF_DB_PORT=%d", 5432+offset),
    fmt.Sprintf("PUBGOLF_PORT=%d", 5000+offset),
)
```

Add a `--port-offset` flag to `run` and `test` commands that sets
`PUBGOLF_PORT_OFFSET` in the process environment before `worktreePortOffset()`
reads it.

### 3. Volume Isolation — `run.go`

Inject the worktree-specific data path:

```go
env = append(env,
    fmt.Sprintf("PUBGOLF_DB_HOST_DATA_PATH=%s",
        filepath.Join(projectRoot, worktreeDataDir("data/postgres"))),
)
```

### 4. Coverage Dir Isolation — `test.go`

Update the coverage output directory to use the worktree slug:

```go
coverageDir := filepath.Join(projectRoot, worktreeDataDir("data/go-test-coverage"))
```

### 5. `PUBGOLF_SHARED_DB_URL` with Offset Port — `test.go`

When running in shared-postgres mode, the shared DB URL must reflect the
offset port. After getting env from the provider, override the URL:

```go
if localFlag {
    offset := worktreePortOffset()
    // Build shared DB URL with offset port
    sharedURL := fmt.Sprintf("postgres://pubgolf_dev:pubgolf_dev@localhost:%d/pubgolf_dev?sslmode=disable",
        5432+offset)
    env = append(env, "PUBGOLF_SHARED_DB_URL="+sharedURL)
}
```

### 6. `PUBGOLF_WORKTREE_SLUG` Injection — `test.go`

Inject the slug so `dbtest.NewConn()` can namespace databases:

```go
if slug, _ := worktreeSlug(); slug != "" {
    env = append(env, "PUBGOLF_WORKTREE_SLUG="+slug)
}
```

### 7. `dbtest.go` — Shared Mode Namespace

In `api/internal/lib/dbtest/dbtest.go`, modify the database name construction
for shared mode:

```go
func sharedDBName(namespace string) string {
    base := "pubgolf_api_dbtest_" + namespace
    if slug := os.Getenv("PUBGOLF_WORKTREE_SLUG"); slug != "" {
        base += "_" + strings.ReplaceAll(slug, "-", "_")
    }
    return base
}
```

### 8. Pre-Flight Health Check — `test.go`

Before running tests, perform a quick infrastructure check:

```go
func preflight(ctx context.Context, offset int) error {
    dbPort := 5432 + offset
    conn, err := net.DialTimeout("tcp",
        fmt.Sprintf("localhost:%d", dbPort),
        2*time.Second)
    if err != nil {
        return fmt.Errorf("DB container is not running on port %d.\n"+
            "  Run 'pubgolf-devctrl run bg' to start it.\n"+
            "  This is an infrastructure issue, not a code problem.", dbPort)
    }
    conn.Close()
    return nil
}
```

Call this before invoking the test command. If it fails, the error will be
classified as infrastructure (exit code 2) by the error handling from plan 05.

### 9. Post-Start Reminder — `run.go`

After `dockerRun` or `goRun` starts successfully from a worktree:

```go
if slug, _ := worktreeSlug(); slug != "" {
    log.Printf("[pubgolf-devctrl] Started services for worktree %q (DB: port %d).\n"+
        "  Run 'pubgolf-devctrl stop' before removing this worktree.", slug, 5432+offset)
}
```

### 10. SIGKILL → SIGINT Fix

The `ExecRunner.Start()` from plan 01 already sends SIGINT. Verify that the
existing SIGKILL behavior in `run.go` is fully replaced — no direct
`syscall.Kill` calls should remain.

### 11. Verify

- Start services from two worktrees concurrently:
  - `cd .worktrees/worktree-devtools && pubgolf-devctrl run bg`
  - `cd .worktrees/test-wt && pubgolf-devctrl run bg`
  - Verify different ports via `docker ps`
- `pubgolf-devctrl stop` from each worktree stops only that worktree's services
- `go test ./tools/...` — passes
- `pubgolf-devctrl check:go` — passes

## Acceptance Criteria

- [ ] Docker project names are worktree-specific
- [ ] Ports are offset per worktree
- [ ] Volume paths are worktree-specific
- [ ] Coverage dir is worktree-specific
- [ ] `PUBGOLF_SHARED_DB_URL` reflects offset port
- [ ] `PUBGOLF_WORKTREE_SLUG` injected into test env
- [ ] `dbtest.go` namespaces shared-mode databases by slug
- [ ] Pre-flight health check before tests
- [ ] Post-start reminder from worktrees
- [ ] No direct `syscall.Kill` / SIGKILL in run.go or test.go
- [ ] `pubgolf-devctrl stop` from worktree stops only that worktree
- [ ] `pubgolf-devctrl check:go` passes
