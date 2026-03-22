# Plan 04: Worktree Exclusions (Lint/Buf/Watchers)

**Depends on:** nothing
**Branch:** `devctrl/04-worktree-exclusions`
**PR scope:** `tools/lib/cmd/check.go`, `tools/lib/cmd/root.go`, `buf.yaml`

## Objective

Prevent linting, buf, and file watchers on main from picking up files inside
`.worktrees/`. This is the most immediately-felt bug: lint failures on main
caused by in-flight worktree changes.

## Steps

### 1. golangci-lint — `check.go`

Add `--exclude-dirs` to the lint invocation in `checkGo()`:

```go
lint := exec.CommandContext(ctx, "golangci-lint", "run",
    "--exclude-dirs", `\.worktrees`,
    "./api/...", "./tools/...",
)
```

### 2. buf — `buf.yaml`

Read the current `buf.yaml` and add `.worktrees` to the `excludes` list. If
no `excludes` key exists, add one. This affects both `buf lint` and
`buf generate`.

### 3. File Watchers — `root.go`

Add a centralized ignore list and apply it in the `watch()` function:

```go
var ignoredDirPatterns = []string{
    ".worktrees",
    "node_modules",
    "vendor",
    "data",
    ".git",
}

func watch(dir, label string, callback func(watcher.Event)) {
    w := watcher.New()
    w.SetMaxEvents(1)

    // Apply ignore patterns
    for _, pattern := range ignoredDirPatterns {
        w.Ignore(pattern)  // or use w.AddFilterHook if needed
    }
    // ... rest unchanged
}
```

Note: check the `radovskyb/watcher` API for the correct ignore mechanism. It
may be `w.Ignore()` or a filter hook. Read the library docs and use whichever
is appropriate.

### 4. Verify

- Create a temporary worktree with a Go file that has lint issues:
  `git worktree add .worktrees/test-exclude -b test-exclude`
  Add a file with an unused variable to `.worktrees/test-exclude/api/`.
- Run `pubgolf-devctrl check:go` from the project root — should NOT report
  the worktree file.
- Clean up: `git worktree remove .worktrees/test-exclude`
- `pubgolf-devctrl check:proto` — should still work normally.
- `pubgolf-devctrl check:go` — no lint issues in the actual code.

## Acceptance Criteria

- [ ] `golangci-lint` excludes `.worktrees/` directory
- [ ] `buf.yaml` excludes `.worktrees/`
- [ ] File watchers ignore `.worktrees/`, `node_modules/`, `vendor/`, `data/`, `.git/`
- [ ] Lint on main is unaffected by worktree contents
- [ ] `pubgolf-devctrl check:go` passes
- [ ] `pubgolf-devctrl check:proto` passes
