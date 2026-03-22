# Plan 02: Worktree Identity + Port Offset

**Depends on:** nothing
**Branch:** `devctrl/02-worktree-identity`
**PR scope:** `tools/lib/cmd/worktree.go`, `tools/lib/cmd/worktree_test.go`

## Objective

Implement worktree detection, slug normalization, and port offset calculation.
These are pure functions with no side effects — fully testable without any
infrastructure.

## Steps

### 1. Create `tools/lib/cmd/worktree.go`

#### `worktreeSlug() (string, error)`

- Shell out to `git rev-parse --show-toplevel` and `git rev-parse --git-common-dir`.
- Resolve both to absolute paths via `filepath.Abs` before comparison.
- If `commonDir` is not a prefix of `topLevel`, we're in a worktree.
  Return `normalizeSlug(filepath.Base(topLevel))`.
- If we are in the main working tree, return `("", nil)`.
- If git commands fail, return `("", err)` — callers must handle this as an
  error, not a fallback to main-tree behavior.

#### `normalizeSlug(raw string) string`

- Lowercase the input.
- Replace non-alphanumeric characters with hyphens.
- Collapse consecutive hyphens.
- Trim leading/trailing hyphens.
- If result exceeds 20 characters: truncate to 20 and append `-` plus the
  first 6 hex characters of the FNV-32a hash of the original (pre-truncation)
  string. Total max: 27 characters.

Examples:
- `fix-auth` → `fix-auth`
- `feature/minio-s3-support` → `minio-s3-support` (filepath.Base already strips)
- `issue-1234-some-very-long-description` → `issue-1234-some-very-a3f2b1`

#### `worktreePortOffset() int`

- Call `worktreeSlug()`. If slug is empty (main tree), return 0.
- Check `PUBGOLF_PORT_OFFSET` env var. If set and valid integer in [1, 500),
  return that value.
- Otherwise: FNV-32a hash of slug, `% 500 + 1`.

#### Helper: `worktreeDockerProject() string`

- Returns `"pubgolf"` for main tree, `"pubgolf-" + slug` for worktrees.

#### Helper: `worktreeDataDir(base string) string`

- Returns `base` for main tree, `base + "-" + slug` for worktrees.
- Used for: `data/postgres`, `data/go-test-coverage`.

### 2. Create `tools/lib/cmd/worktree_test.go`

#### Slug normalization tests

```go
func TestNormalizeSlug(t *testing.T) {
    tests := []struct{ in, want string }{
        {"fix-auth", "fix-auth"},
        {"UPPER-Case", "upper-case"},
        {"special!@#chars", "special-chars"},
        {"a--b--c", "a-b-c"},
        {"-leading-trailing-", "leading-trailing"},
        {"issue-1234-some-very-long-description-of-the-bug", /* truncated + hash */},
    }
    // ...
}
```

#### Port offset tests

- Determinism: same slug always returns same offset.
- Range: offset is in [1, 500].
- Env var override: set `PUBGOLF_PORT_OFFSET=42`, verify it takes precedence.
- Env var validation: out-of-range values fall through to hash.
- Main tree: returns 0.

#### Docker project name tests

- Main tree: `"pubgolf"`.
- Worktree: `"pubgolf-fix-auth"`.

#### Data dir tests

- Main tree: `"data/postgres"`.
- Worktree: `"data/postgres-fix-auth"`.

### 3. Verify

- `go test ./tools/...` — all new tests pass
- `pubgolf-devctrl check:go` — no lint issues

## Acceptance Criteria

- [ ] `worktreeSlug()` returns `(string, error)`, distinguishes main/worktree/error
- [ ] `normalizeSlug` handles all edge cases with tests
- [ ] `worktreePortOffset` supports `PUBGOLF_PORT_OFFSET` env var override
- [ ] All helper functions tested
- [ ] No dependencies on Runner, EnvProvider, or any other plan's code
- [ ] `go test ./tools/...` passes
- [ ] `pubgolf-devctrl check:go` passes
