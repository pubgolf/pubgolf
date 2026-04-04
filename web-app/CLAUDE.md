# web-app/ — SvelteKit Frontend

SvelteKit SPA using adapter-static, Svelte 5 (compat mode), Skeleton UI 1.x, and Tailwind CSS.
Connects to the Go API server via Connect-RPC (web transport).

## Dev Workflow

All commands go through `pubgolf-devctrl` (run from project root, not `web-app/`).

| Task | Command |
|------|---------|
| Dev server (HMR, port 5173) | `pubgolf-devctrl run web` |
| Full-stack (DB + preview + API) | `pubgolf-devctrl run` |
| Production build | `pubgolf-devctrl build web` |
| Install npm deps | `pubgolf-devctrl install web` |
| Unit tests (vitest) | `pubgolf-devctrl test web` |
| E2E tests (Playwright) | `pubgolf-devctrl test e2e web` |
| Lint + type-check | `pubgolf-devctrl check web` |

**Full-stack mode** (`pubgolf-devctrl run` with no subcommand): starts the DB via
Docker, starts `vite dev` (with HMR), and starts the API server reverse-proxying
to the dev server. Access the app through the API server port. Can be backgrounded:

    pubgolf-devctrl run &  # save the PID from $!
    kill -INT <pid>        # clean shutdown (stops all child processes)

`pubgolf-devctrl stop` only tears down Docker — use SIGINT on the `run` process
to stop everything.

**Fresh worktree**: Run `pubgolf-devctrl migrate up` before the first `run` —
the DB starts empty and the API server will panic without tables.

**Do not use `npm run` directly** — devctrl is the single source of truth.
The CI pipeline also runs through devctrl (`go run ./tools/cmd/pubgolf-devctrl`).

**Dependency conflicts**: Use `overrides` in `package.json` to resolve peer dep
conflicts (e.g., Skeleton v1 + Svelte 5). Do not use `--legacy-peer-deps` —
CI runs `npm ci` which doesn't honor it.

## Key Directories

- `src/routes/` — SvelteKit file-based routing
- `src/lib/components/` — shared Svelte components
- `src/lib/helpers/` — pure utility functions (with unit tests)
- `src/lib/proto/` — generated Connect-RPC client types (do not edit)
- `src/lib/rpc/` — RPC client setup
- `e2e-tests/` — Playwright e2e test specs
- `static/` — static assets copied to build output

## Proto / RPC Client

The TypeScript proto types in `src/lib/proto/` are generated from `proto/` by
`pubgolf-devctrl generate proto`. Never edit these files directly — edit the
`.proto` source and regenerate.

The RPC client in `src/lib/rpc/client.ts` configures Connect-RPC web transport.

## Assets

Static assets are served from `assets.pubgolf.co` (Cloudflare R2) in production.
The base path is configured via `SVELTE_ASSETS_PATH` in `svelte.config.js`.
Immutable assets (`_app/immutable/`) are uploaded to R2 during CI.

## Testing

- **Unit tests**: vitest, run via `pubgolf-devctrl test web`. Test files live
  alongside their source in `src/lib/helpers/` (e.g. `formatters.test.ts`).
- **E2E tests**: Playwright, run via `pubgolf-devctrl test e2e web`. Specs
  live in `e2e-tests/`. Playwright config handles build + preview automatically.
