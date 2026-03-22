# Plan 03: Robustness Fixes (Non-Runner)

**Depends on:** nothing
**Branch:** `devctrl/03-robustness-fixes`
**PR scope:** Targeted fixes to `root.go`, `stop.go`, `migrate.go`, `config.go`

## Objective

Fix the confirmed brittleness issues that do NOT depend on the Runner
interface. These are small, independent fixes that reduce risk of the tool
breaking in normal use.

## Steps

### 1. `stopCmd` Description — `stop.go`

**Problem:** `config.CLIName` is empty at package init time.

**Fix:** Hardcode the description:
```go
var stopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Stop all background processes started with `pubgolf-devctrl run ...`",
    // ...
}
```

### 2. `migrateCreate` Stderr Parsing — `migrate.go`

**Problem:** `outLines[:len(outLines)-1]` can panic on empty/unexpected output.

**Fix:** Replace the parsing block:
```go
outLines := strings.Split(migratorContent.String(), "\n")
for _, name := range outLines {
    name = strings.TrimSpace(name)
    if name == "" || !strings.Contains(name, migrationDirectory) {
        continue
    }
    // ... open and write boilerplate
}
```

If no valid file paths found, log a warning instead of silently continuing.

### 3. Project Root Resolution — `root.go` + `config.go`

**Problem:** All relative paths assume CWD is project root.

**Fix:** Add a `resolveProjectRoot() (string, error)` function:
- Check CWD for `go.mod`. If found, return CWD.
- Walk up parent directories looking for `go.mod`.
- If found, return that directory path.
- If not found, return error.

Store result as `var projectRoot string` (package-level). Set in `Execute()`
before `rootCmd.Execute()`.

Log: `[pubgolf-devctrl] Resolved project root: /path/to/root`

**Do NOT call `os.Chdir`.** Store the path and use `filepath.Join(projectRoot, ...)`
in command construction. This is safer with concurrent goroutines.

For this PR, update `checkVersion()` to use `filepath.Join(projectRoot, "tools")`
for the `dirhash.HashDir` call. Other command files will be updated to use
`projectRoot` as part of the Runner threading (plan 01) — don't duplicate that
work here.

### 4. `beginShutdown` Close Safety — `root.go`, `run.go`, `test.go`

**Problem:** `close(beginShutdown)` can panic if called twice.

**Fix:**
```go
var shutdownOnce sync.Once

func triggerShutdown() {
    shutdownOnce.Do(func() { close(shuttingDown) })
}
```

Replace all `close(beginShutdown)` calls with `triggerShutdown()`.
Update the signal handler goroutine in `PersistentPreRun`:
```go
go func() {
    <-beginShutdown
    triggerShutdown()
}()
```

### 5. Verify

- `pubgolf-devctrl check:go` — no lint issues
- `pubgolf-devctrl migrate create test_fix` — verify boilerplate is written
  (then delete the test migration files)
- Manual: run devctrl from a subdirectory, verify it resolves to project root

## Acceptance Criteria

- [ ] `stopCmd.Short` displays correctly
- [ ] `migrateCreate` handles empty/unexpected stderr gracefully
- [ ] `projectRoot` resolved and used for `checkVersion` dirhash
- [ ] `triggerShutdown()` replaces all `close(beginShutdown)` calls
- [ ] `pubgolf-devctrl check:go` passes
- [ ] No functional changes to command behavior (pure bug fixes)
