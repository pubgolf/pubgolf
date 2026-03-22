# Plan 08: New Commands (status, clean, doctor, stop enhancements)

**Depends on:** 02-worktree-identity
**Branch:** `devctrl/08-new-commands`
**PR scope:** `tools/lib/cmd/status.go`, `tools/lib/cmd/clean.go`,
`tools/lib/cmd/doctor.go`, modifications to `tools/lib/cmd/stop.go`

## Objective

Add agent-oriented commands for introspection, cleanup, environment validation,
and enhanced stop behavior. These commands make agents self-sufficient by
providing the diagnostic information they need without shelling out to Docker/git
directly.

## Prerequisites

Plan 02 (worktree identity) must be merged. These commands use `worktreeSlug()`,
`worktreePortOffset()`, `worktreeDockerProject()`, and `worktreeDataDir()`.

Note: These commands use `exec.CommandContext` directly (not via Runner) for
Docker/git introspection calls. Once plan 01 (Runner) is also merged, they can
be updated to use Runner, but this is not blocking — the commands are read-only
and don't benefit from dry-run recording.

## Steps

### 1. `devctrl status` — `status.go`

```go
var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show worktree identity, port assignments, and service status",
    Run: func(cmd *cobra.Command, _ []string) {
        allFlag, _ := cmd.Flags().GetBool("all")
        if allFlag {
            printAllWorktreeStatus()
        } else {
            printCurrentStatus()
        }
    },
}
```

`printCurrentStatus()`:
- Print worktree slug (or "(main working tree)")
- Port offset and resolved ports (DB, API)
- Docker project name
- Data volume path
- Docker container status (query `docker ps --filter name=<project>`)

`printAllWorktreeStatus()`:
- Parse `git worktree list --porcelain` to enumerate worktrees
- For each, compute slug and offset
- Check Docker status for each project name
- Print as aligned table

### 2. `devctrl clean` — `clean.go`

Defaults to dry-run. Must pass `--force` to actually delete.

```go
var cleanCmd = &cobra.Command{
    Use:   "clean",
    Short: "Remove orphaned worktree resources (Docker containers, data dirs)",
    Run: func(cmd *cobra.Command, _ []string) {
        force, _ := cmd.Flags().GetBool("force")
        cleanOrphans(cmd.Context(), force)
    },
}
```

**Two-pass algorithm:**

Pass 1 — Docker (skip if daemon unavailable):
- List Docker Compose projects matching `pubgolf-*` via
  `docker compose ls --all --format json`.
- Get active worktree slugs from `git worktree list`.
- Projects whose slug suffix doesn't match an active worktree are orphaned.
- If `--force`: `docker compose --project-name <name> down --volumes`.
- Otherwise: print what would be removed.

Pass 2 — Filesystem (always runs):
- Scan `data/postgres-*` and `data/go-test-coverage-*` in project root.
- Extract slug suffix, compare against active worktree slugs.
- If `--force`: `os.RemoveAll`.
- Otherwise: print what would be removed (with size via `du -sh` equivalent).

If Docker is unavailable:
```
[pubgolf-devctrl] WARNING: Docker daemon is not available. Skipping container cleanup.
  Start Docker and re-run to clean up Docker resources.
```

### 3. `devctrl doctor` — `doctor.go`

Comprehensive environment check. All checks are read-only.

```go
var doctorCmd = &cobra.Command{
    Use:   "doctor",
    Short: "Check that all development tools are installed and configured",
    Run: func(cmd *cobra.Command, _ []string) {
        runDoctorChecks(cmd.Context())
    },
}
```

Checks:
- `go version` — present, version ≥ 1.26
- `docker info` — Docker daemon running
- `doppler --version` — installed (note: does not verify access)
- `golangci-lint --version` — installed
- `buf --version` — installed
- `sqlc version` — installed
- `mockery --version` — installed
- `node --version` — installed
- `npm --version` — installed
- Project root detected

Output format:
```
[pubgolf-devctrl] Environment check:
  Go:             1.26.0    ✓
  Docker:         running   ✓
  Doppler:        1.2.3     ✓
  golangci-lint:  2.11.3    ✓
  buf:            1.30.0    ✓
  sqlc:           1.24.0    ✓
  mockery:        2.42.2    ✓
  Node.js:        22.x      ✓
  npm:            10.x      ✓
  Project root:   /path     ✓
```

On failure:
```
  Docker:         not running  ✗  (run 'open -a Docker' or 'dockerd')
```

### 4. `devctrl stop` Enhancements — `stop.go`

Add two flags:

#### `--all`

Stops services across all worktree projects:
- Enumerate worktrees via `git worktree list`
- For each, compute Docker project name
- Run `docker-compose --project-name <name> down` for each

#### `--remove-data`

After stopping, remove this worktree's data directories:
- `data/postgres-{slug}`
- `data/go-test-coverage-{slug}`

For explicit teardown when an agent is done with a worktree.

### 5. Register Commands

In `init()` blocks or in the respective files:
- `rootCmd.AddCommand(statusCmd)`
- `rootCmd.AddCommand(cleanCmd)`
- `rootCmd.AddCommand(doctorCmd)`
- Add `--all` and `--remove-data` flags to `stopCmd`

### 6. Verify

- `pubgolf-devctrl status` — prints current worktree info
- `pubgolf-devctrl status --all` — lists all worktrees
- `pubgolf-devctrl clean` — dry-run, shows what would be cleaned
- `pubgolf-devctrl doctor` — shows tool versions
- `pubgolf-devctrl stop --all` — stops all projects
- `pubgolf-devctrl check:go` — no lint issues

## Acceptance Criteria

- [ ] `devctrl status` prints current worktree identity and ports
- [ ] `devctrl status --all` lists all worktrees with Docker status
- [ ] `devctrl clean` defaults to dry-run
- [ ] `devctrl clean --force` removes orphaned resources
- [ ] `devctrl clean` handles Docker unavailability gracefully
- [ ] `devctrl doctor` checks all required tools
- [ ] `devctrl stop --all` stops all worktree projects
- [ ] `devctrl stop --remove-data` cleans up data dirs
- [ ] `pubgolf-devctrl check:go` passes
