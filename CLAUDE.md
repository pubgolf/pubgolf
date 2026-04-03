# PubGolf — Claude Instructions

## Project Overview
Go (Connect-RPC gRPC) backend + SvelteKit frontend + proto-first API.
Module: `github.com/pubgolf/pubgolf` | Go 1.26

## Task Runner: pubgolf-devctrl
`pubgolf-devctrl` is the task runner. The following tools are hard-blocked because
devctrl handles them: `sqlc`, `mockery`, `enumer`, `ifacemaker`, `migrate`, `go run`, `go generate`.

The following tools require approval (prefer the devctrl equivalent):
`buf`, `golangci-lint`, `go build`, `go mod`, `npm run`.

**Syntax**: devctrl uses space-separated subcommands — e.g. `pubgolf-devctrl check go`,
NOT `pubgolf-devctrl check:go`. The colon notation in the table below (`check:*`) is
shorthand for permission grouping, not the actual CLI syntax.

| Task | Command (shorthand) | Status |
|------|---------------------|--------|
| Code gen (proto, SQL, mocks, enums) | `pubgolf-devctrl generate *` | auto-approved |
| DB migrations | `pubgolf-devctrl migrate *` | auto-approved |
| Run full-stack or individual servers | `pubgolf-devctrl run [api|web|bg|api-db]` | auto-approved |
| Build targets | `pubgolf-devctrl build [web]` | auto-approved |
| Run tests (Go + web + e2e) | `pubgolf-devctrl test [web|e2e [api|web]]` | auto-approved |
| Lint/check all packages | `pubgolf-devctrl check go` | auto-approved |
| Lint/check proto files | `pubgolf-devctrl check proto` | auto-approved |
| Check web app (lint + type-check) | `pubgolf-devctrl check web` | auto-approved |
| Install dev deps | `pubgolf-devctrl install *` | auto-approved |
| Stop servers | `pubgolf-devctrl stop *` | auto-approved |
| Inspect worktree status | `pubgolf-devctrl status` | auto-approved |
| Inspect all worktrees | `pubgolf-devctrl status --all` | auto-approved |
| Check dev environment | `pubgolf-devctrl doctor` | auto-approved |
| Preview cleanup | `pubgolf-devctrl clean` | auto-approved |
| Dry-run any command | `pubgolf-devctrl --dry-run <cmd>` | auto-approved |
| Update devctrl itself | `pubgolf-devctrl update *` | requires approval |
| Force cleanup | `pubgolf-devctrl clean --force` | requires approval |

**Single-package Go tests**: `pubgolf-devctrl test` always runs the full suite with
Doppler secrets. For iterating on a single package (no DB/secrets needed), use
`go test ./api/internal/lib/your/pkg/...` directly — this is pre-approved.

## Code Exploration
Use dedicated tools: **Glob**, **Grep**, **Read** instead of bash find/grep/cat.
(Note: if you are not Claude Code, use standard bash equivalents like ls, grep -r, cat.)

`find` is not pre-approved. `xargs`, `eval`, `source` are blocked globally.

## Sensitive Configuration Files
The following files are treated as executable code — changes to them can affect
what runs during code generation or builds and should be reviewed carefully:
- `buf.gen.dev.yaml` / `buf.gen.ci.yaml` — buf plugin configuration
- `api/internal/db/sqlc.yaml` — sqlc codegen configuration
- `.git/hooks/` — do NOT edit git hooks directly

`go install` fetches and executes external packages. Only use with pinned versions
from known modules, not `@latest`.

## Scripting
Do NOT write ad-hoc `.sh`/`.py`/`.js` files to run one-off tasks — script execution
by filename is blocked. For persistent tooling, add commands to
`tools/cmd/pubgolf-devctrl/`. For one-off operations, write Go test functions or
use devctrl subcommands.

## Key Directories
- `api/` — Go API server (Connect-RPC)
- `web-app/` — SvelteKit frontend
- `tools/` — pubgolf-devctrl source
- `proto/` — API definitions (source of truth; never edit generated files directly)

## Proto Workflow
1. Edit `.proto` files
2. `pubgolf-devctrl generate proto` (regenerates all client libs)
3. `pubgolf-devctrl check proto` (runs buf lint + format diff)
4. Commit generated files

Always run generate before check — check will fail on stale generated output.

## Definition of Done
Before declaring any task complete, run the relevant verification commands.
Do not use "equivalent" alternatives — use the canonical command.

**Before submitting or non-trivially updating a PR**, run `/simplify` to review
changed code for reuse, quality, and efficiency.

| Change type | Verification command |
|-------------|---------------------|
| Go code changes | `pubgolf-devctrl check go` then `go test ./affected/pkg/...` |
| Proto file changes | `pubgolf-devctrl generate proto` then `pubgolf-devctrl check proto` |
| Codegen-affecting changes | `pubgolf-devctrl generate *` before running tests |
| Web app changes | `pubgolf-devctrl check web` |

**Web verification is always pre-approved via devctrl.** Do not use `npm run check` or
`npm run ci:lint` directly — `pubgolf-devctrl check web` runs both automatically.

## Autonomous and Subagent Tasks
When running as a subagent or background task, approval prompts cannot be fulfilled.
Design tasks to use only pre-approved devctrl commands and `go test`/`go vet`.

If a task requires `npm run *`, `go build`, `buf`, or other ask-gated commands,
surface this as a limitation before starting — provide the commands for the user
to run manually rather than blocking mid-task.

`pubgolf-devctrl update` will trigger an approval prompt mid-task. Do not invoke it
as part of an automated sequence unless it is necessary specifically for tool development.

## Exit Codes
- Exit 1: test/lint failure (debug your code).
- Exit 2: infrastructure failure (check services, not your code).

## Worktree Port Collisions
If you encounter "address already in use", set `PUBGOLF_PORT_OFFSET` to a
number between 1 and 499 and retry:

    export PUBGOLF_PORT_OFFSET=42

## Testing
Go tests use [Testify](https://github.com/stretchr/testify). Use `require` for
setup/preconditions (hard-fail) and `assert` for result assertions (soft-fail).
Do not use raw `t.Errorf`/`t.Fatalf` — always prefer `assert.*` / `require.*`.

## Git
`git push` always requires explicit approval.
Do not edit `.git/hooks/` or run `git config --global`.

## Network
Prefer `WebFetch` / `WebSearch` tools for read-only access.
`curl`, `wget` require approval. `npx` is blocked globally.
