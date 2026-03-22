# Plan 05: Error Handling Refactor

**Depends on:** 01-runner-interface
**Branch:** `devctrl/05-error-handling`
**PR scope:** `tools/lib/cmd/errors.go`, modifications to all command files

## Objective

Replace `guard`/`log.Fatalf` with classified error returns and distinct exit
codes. This is critical for agent autonomy — agents need to distinguish "your
code is broken" (exit 1) from "infrastructure is broken" (exit 2).

## Prerequisites

Plan 01 must be merged first. The Runner refactor converts all extracted
functions to accept `runner Runner` and return `error`. This plan builds on
that by making the error returns meaningful.

## Steps

### 1. Create `tools/lib/cmd/errors.go`

```go
const (
    ExitCodeTestFailure   = 1
    ExitCodeInfraFailure  = 2
)

// isInfraError checks stderr content for known infrastructure failure patterns.
func isInfraError(err error) bool {
    msg := err.Error()
    patterns := []string{
        "address already in use",
        "Cannot connect to the Docker daemon",
        "connection refused",
        "permission denied",
        "no such host",
        "container exited with code 137",  // OOM killed
    }
    for _, p := range patterns {
        if strings.Contains(strings.ToLower(msg), strings.ToLower(p)) {
            return true
        }
    }
    return false
}

// classifyAndExit logs the error with classification and exits.
func classifyAndExit(err error) {
    if err == nil {
        return
    }
    if isInfraError(err) {
        log.Printf("[pubgolf-devctrl] ERROR: Infrastructure failure\n  %s\n  This is not a code issue.", err)
        os.Exit(ExitCodeInfraFailure)
    }
    log.Printf("[pubgolf-devctrl] ERROR: %s", err)
    os.Exit(ExitCodeTestFailure)
}
```

### 2. Refactor `runPar` — `root.go`

Replace the current goroutine-with-guard pattern:

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

### 3. Replace `guard` calls with error returns

For each command file, the functions should now return `error` (plan 01 already
does this). What changes here is the cobra `Run` closures — instead of calling
`guard(err, ...)`, they call `classifyAndExit(err)`:

```go
var checkGoCmd = &cobra.Command{
    Use:   "go",
    Short: "Run golangci-lint on Go code",
    Run: func(cmd *cobra.Command, _ []string) {
        classifyAndExit(checkGo(cmd.Context(), runner))
    },
}
```

### 4. Port collision error message

In the error wrapping for Docker/server start commands, wrap port-bind errors
with actionable context:

```go
if strings.Contains(err.Error(), "address already in use") {
    return fmt.Errorf("port %d is already in use (%s)\n"+
        "  Worktree %q was assigned port offset %d (base %d + %d = %d).\n"+
        "  To resolve:\n"+
        "    1. Run: pubgolf-devctrl stop          (in the conflicting worktree)\n"+
        "    2. Or:  export PUBGOLF_PORT_OFFSET=%d && pubgolf-devctrl run bg\n"+
        "    3. Or:  pubgolf-devctrl clean --force  (from main)",
        port, envVarName, slug, offset, basePort, offset, port, offset+1)
}
```

### 5. Dry-run sequence numbers

Now that `runPar` is refactored, add sequence numbering to the dry-run output.
The DryRunner tracks a counter and annotates each recorded command with its
position and whether it was part of a parallel group.

### 6. Verify

- `pubgolf-devctrl check:go` — still works, errors display correctly
- `go test ./tools/...` — error classification tests pass
- Manually test an infrastructure failure (e.g., stop Docker, run a command)
  — verify exit code 2 and clear error message

## Acceptance Criteria

- [ ] `guard` function removed entirely (or deprecated with a TODO)
- [ ] All cobra `Run` closures use `classifyAndExit`
- [ ] `runPar` collects errors instead of calling `log.Fatalf`
- [ ] `isInfraError` classifies known infrastructure patterns
- [ ] Exit code 1 for code/test failures, exit code 2 for infrastructure
- [ ] Port collision error messages include remediation steps
- [ ] `pubgolf-devctrl check:go` passes
- [ ] `go test ./tools/...` passes
