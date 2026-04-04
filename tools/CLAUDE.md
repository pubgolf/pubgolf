# tools/ — DevCtrl Conventions

All devctrl source lives under `tools/lib/cmd/`. Each file maps to a top-level
command group (`check.go` → `devctrl check`, `generate.go` → `devctrl generate`, etc.).
The entry point in `tools/cmd/pubgolf-devctrl/main.go` just calls `cmd.Execute()`.

## Testing Source Changes

`pubgolf-devctrl` on PATH is a compiled binary — editing source files in a
worktree has no effect until you recompile with `pubgolf-devctrl update`.
The update output includes a content hash; a changed hash confirms new source
was compiled.

## Adding a New Command

1. Create a file named after the command group (e.g. `deploy.go`).
2. Register in `init()` — add subcommands to the parent, then the parent to `rootCmd`.
3. Cobra handler should be a thin wrapper that calls a testable function:

```go
func init() {
    deployCmd.AddCommand(deployAPICmd)
    rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy all services",
    Run: func(cmd *cobra.Command, _ []string) {
        classifyAndExit(deployAll(cmd.Context(), runner))
    },
}

// deployAll is the testable handler — accepts Runner, returns error.
func deployAll(ctx context.Context, r Runner) error {
    err := r.Run(ctx, Cmd{Name: "kubectl", Args: []string{"apply", "-f", "..."}})
    if err != nil {
        return fmtErr(err, "run kubectl apply cmd")
    }

    return nil
}
```

## Error Handling

Two functions, always used together:
- `fmtErr(err, "context message")` — wraps with context, nil-safe.
- `classifyAndExit(err)` — logs and exits with code 1 (test/lint) or 2 (infra).

Handlers **return errors**. Only the cobra `Run` wrapper calls `classifyAndExit`:

```go
// In handler (returns error):
return fmtErr(err, "run buf lint cmd")

// In cobra Run (exits):
classifyAndExit(handlerFunc(cmd.Context(), runner))
```

Never call `log.Fatal`, `os.Exit`, or `panic` directly from handler functions.

## Runner & Process

All subprocess execution goes through the `Runner` interface:
- `Runner.Run(ctx, Cmd)` — execute and wait.
- `Runner.Start(ctx, Cmd)` — start a long-running process, returns `Process`.
- `Process.Wait()` / `Process.Stop()` — lifecycle management.

Production: `ExecRunner`. Tests/dry-run: `DryRunner` (records commands, injects errors).

Handler functions always accept `Runner` as a parameter — never use the package-level
`runner` directly in a handler. This enables testing with `DryRunner`.

```go
// Good: accepts Runner parameter
func checkGo(ctx context.Context, r Runner) error { ... }

// Bad: uses package global
func checkGo(ctx context.Context) error { runner.Run(...) }
```

## EnvProvider

Commands that need secrets accept `EnvProvider` as a parameter:
- `DopplerProvider` — fetches from Doppler CLI.
- `EnvFileProvider` — reads `.env.{config}` files.
- `AutoProvider` — tries Doppler, falls back to env files.

Same rule as Runner: pass as parameter, don't use the package global in handlers.

## Parallel Execution

Use `runPar` to run independent tasks concurrently. It collects all errors via
`errors.Join` (not just the first):

```go
classifyAndExit(runPar(cmd.Context(), runner, checkGo, checkWeb, checkProto))
```

Each function passed to `runPar` must have the signature
`func(context.Context, Runner) error`.

## Testing

- Use `DryRunner` to verify which commands a handler invokes:
  ```go
  dr := &DryRunner{}
  err := deployAll(context.Background(), dr)
  require.NoError(t, err)
  require.Len(t, dr.Recorded, 1)
  assert.Equal(t, "kubectl", dr.Recorded[0].Name)
  ```
- Use `DryRunner.ErrorFor` to simulate failures by command name.
- Use `require.*` for preconditions, `assert.*` for result checks (testify).
- All tests must call `t.Parallel()`.

## Naming

- Files: one per command group, named to match the cobra `Use` field.
- Handler functions: `verbNoun` (e.g. `checkGo`, `generateProto`, `dockerRun`).
- Cobra vars: `{group}Cmd`, `{group}{Sub}Cmd` (e.g. `checkCmd`, `checkGoCmd`).

## Function Signatures

**Rule of thumb: 4 parameters max** (not counting `ctx context.Context`). When a
function exceeds this, simplify in order of preference:

1. **Eliminate redundant params** — if every caller passes the same value (e.g.
   `config.ServerBinName`), read it from the package-level `config` directly.
   Config is static and doesn't need to be injected for testability.
2. **Group into an opts struct** — last resort. Structs trade compile-time
   completeness (missing positional param = compiler error) for readability.
   Only worth it when params genuinely vary across callers.

`Runner` and `EnvProvider` stay as direct params — they're swapped in tests
(`DryRunner`, mock providers) so they must be injectable.

## nolint Annotations

- `//nolint:ireturn` — on functions that return an interface by design (Runner, Process, EnvProvider).
- `//nolint:gosec` — on `exec.CommandContext` (args are controlled by devctrl internals).
- `//nolint:errorlint` — on `err.(*exec.ExitError)` casts that extract exit codes.
