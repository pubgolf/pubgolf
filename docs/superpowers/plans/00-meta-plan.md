# Meta Plan: devctrl Worktree Robustness

**Spec:** `docs/superpowers/specs/2026-03-16-devctrl-worktree-robustness-design.md`

## Work Units

Each unit produces one PR. Units are designed so that each can be executed by
a single Claude Code agent in an isolated git worktree.

| ID | Name | Depends On | Parallel Group |
|----|------|-----------|----------------|
| 01 | Runner interface + dry-run | — | A |
| 02 | Worktree identity + port offset | — | A |
| 03 | Robustness fixes (non-Runner) | — | A |
| 04 | Worktree exclusions (lint/buf/watchers) | — | A |
| 05 | Error handling refactor | 01 | B |
| 06 | Env provider decoupling | 01 | B |
| 07 | Docker + test isolation | 01, 02, 06 | C |
| 08 | New commands (status, clean, doctor, stop enhancements) | 02 | B |
| 09 | Test suite + CLAUDE.md | 01, 02, 05 | D |

## Dependency Graph

```
  Group A (all parallel, no deps):
  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │ 01-runner │  │02-worktree│  │03-robust  │  │04-excludes│
  └────┬─────┘  └────┬─────┘  └──────────┘  └──────────┘
       │              │
  Group B (parallel after their deps):
  ┌────▼─────┐  ┌────▼─────┐
  │05-errors │  │06-envprov │  ┌──────────┐
  └────┬─────┘  └────┬─────┘  │08-commands│◄── 02
       │              │        └────┬─────┘
  Group C:            │             │
  ┌───────────────────▼─┐           │
  │ 07-docker-isolation │◄── 02     │
  └─────────┬───────────┘           │
            │                       │
  Group D:  │                       │
  ┌─────────▼───────────────────────▼──┐
  │ 09-tests-and-claudemd              │◄── 01, 02, 05
  └────────────────────────────────────┘
```

## Execution Waves

**Wave 1** (4 agents, fully parallel):
- 01-runner
- 02-worktree
- 03-robustness
- 04-excludes

**Wave 2** (3 agents, after Wave 1 merges):
- 05-errors (needs 01)
- 06-envprovider (needs 01)
- 08-commands (needs 02)

**Wave 3** (1 agent, after Wave 2 merges):
- 07-docker-isolation (needs 01, 02, 06)

**Wave 4** (1 agent, after Wave 3 merges):
- 09-tests-and-claudemd (needs 01, 02, 05)

## Merge Order

Within each wave, PRs can merge in any order. Between waves, all PRs in the
previous wave must be merged before the next wave starts.

In practice, Wave 2 agents can start working as soon as their specific
dependencies merge — they don't need to wait for the entire wave. For example,
05-errors can start as soon as 01-runner merges, even if 02-worktree hasn't
merged yet.

## Notes for Orchestrator

- Each agent gets its own worktree under `.worktrees/plan-{id}/`.
- Agents should use `go test ./tools/...` for validation (no Doppler needed).
- Wave 1 agents work against the current `main` branch.
- Wave 2+ agents should branch from `main` after their dependencies are merged.
- The orchestration script (`run-devctrl-agents.sh`) handles worktree creation
  and agent dispatch.
