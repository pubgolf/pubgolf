#!/usr/bin/env bash
set -euo pipefail

# =============================================================================
# Orchestrator for devctrl worktree robustness implementation
#
# Usage:
#   ./run-devctrl-agents.sh wave1     # Launch Wave 1 agents (no deps)
#   ./run-devctrl-agents.sh wave2     # Launch Wave 2 agents (after Wave 1 merged)
#   ./run-devctrl-agents.sh wave3     # Launch Wave 3 agent  (after Wave 2 merged)
#   ./run-devctrl-agents.sh wave4     # Launch Wave 4 agent  (after Wave 3 merged)
#   ./run-devctrl-agents.sh all       # Launch Wave 1 only (use subsequent waves after merging)
#   ./run-devctrl-agents.sh status    # Show status of all worktrees/branches
#   ./run-devctrl-agents.sh cleanup   # Remove all plan worktrees
# =============================================================================

REPO_ROOT="$(git rev-parse --show-toplevel)"
WORKTREE_DIR="${REPO_ROOT}/.worktrees"
PLANS_DIR="${REPO_ROOT}/docs/superpowers/plans"
SPEC="${REPO_ROOT}/docs/superpowers/specs/2026-03-16-devctrl-worktree-robustness-design.md"
PROMPT_DIR="/tmp/pubgolf-agent-prompts"

mkdir -p "${PROMPT_DIR}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() { echo -e "${BLUE}[orchestrator]${NC} $*"; }
warn() { echo -e "${YELLOW}[orchestrator]${NC} $*"; }
err() { echo -e "${RED}[orchestrator]${NC} $*" >&2; }
ok() { echo -e "${GREEN}[orchestrator]${NC} $*"; }

# Open a new Ghostty tab and run a launcher script
open_ghostty_tab() {
    local launcher="$1"
    local title="$2"

    osascript <<EOF
tell application "Ghostty"
    activate
    set cfg to new surface configuration
    set command of cfg to "bash ${launcher}"

    if (count of windows) is 0 then
        set win to new window with configuration cfg
    else
        set t to new tab in front window with configuration cfg
    end if
end tell
EOF
}

# Launch a single agent in a Ghostty tab with its own worktree
launch_agent() {
    local plan_id="$1"
    local branch="$2"
    local wave_num="$3"
    local plan_file="${PLANS_DIR}/${plan_id}.md"

    if [[ ! -f "${plan_file}" ]]; then
        err "Plan file not found: ${plan_file}"
        return 1
    fi

    local wt_path="${WORKTREE_DIR}/plan-${plan_id}"

    # Create worktree if it doesn't exist
    if [[ ! -d "${wt_path}" ]]; then
        log "Creating worktree: ${wt_path} (branch: ${branch})"
        git worktree add "${wt_path}" -b "${branch}" 2>/dev/null || {
            # Branch may already exist
            git worktree add "${wt_path}" "${branch}" 2>/dev/null || {
                warn "Worktree already exists or branch conflict for ${plan_id}"
            }
        }
    else
        log "Worktree already exists: ${wt_path}"
    fi

    # Write the prompt to a temp file
    local prompt_file="${PROMPT_DIR}/${plan_id}.prompt"
    cat > "${prompt_file}" <<PROMPT
You are implementing a specific plan for the pubgolf-devctrl tool refactor.
This is Wave ${wave_num}, Plan ${plan_id}.

## Your Plan
$(cat "${plan_file}")

## Full Design Spec (for reference)
$(cat "${SPEC}")

## Instructions

### Before you start
1. You are working in a git worktree at: ${wt_path}
2. Your branch is: ${branch}
3. **Read CLAUDE.md first** — it governs all tool usage in this project. Key rules:
   - \`pubgolf-devctrl\` is the canonical task runner. Always prefer devctrl subcommands over raw tools.
   - HARD BLOCKED tools (devctrl handles these): sqlc, mockery, enumer, ifacemaker, migrate, go run, go generate.
   - Approval-required tools (use devctrl instead): buf, golangci-lint, go build, go mod, npm run.
   - Use \`pubgolf-devctrl generate:*\` for codegen, \`pubgolf-devctrl check:go\` for linting, \`pubgolf-devctrl test:*\` for the full test suite.
   - For single-package iteration without Doppler: \`go test ./tools/...\` is pre-approved.

### Implementation
4. Implement the plan step by step, following the acceptance criteria.
5. Do NOT modify files outside the scope listed in the plan.
6. Run \`go test ./tools/...\` after implementation to verify tests pass.
7. Run \`pubgolf-devctrl check:go\` to verify lint passes.

### Finishing up
8. Commit your changes with a descriptive commit message.
9. Push your branch and create a PR:
   - Run: \`git push -u origin ${branch}\`
   - Run: \`gh pr create\` with:
     - Title: "Wave ${wave_num} / Plan ${plan_id}: <short description of what you implemented>"
     - Body should include:
       - A summary section describing what the plan implements
       - A test plan section with checklist items (e.g. \`go test ./tools/...\` passes, \`pubgolf-devctrl check:go\` passes)
       - A merge instructions section: "Merge via GitHub once CI is green and review is approved. This is part of Wave ${wave_num} of the devctrl worktree robustness refactor."
PROMPT

    # Write a launcher script that opens claude interactively in the worktree
    local launcher="${PROMPT_DIR}/${plan_id}.launcher.sh"
    # Resolve paths now — Ghostty tabs launch with --noprofile --norc
    # so nvm/fnm, homebrew, and GOBIN won't be on PATH
    local claude_bin node_dir homebrew_dir gobin_dir
    claude_bin="$(command -v claude)"
    node_dir="$(dirname "$(command -v node)")"
    homebrew_dir="$(dirname "$(command -v go)")"
    gobin_dir="$(dirname "$(command -v pubgolf-devctrl)")"

    cat > "${launcher}" <<LAUNCHER
#!/usr/bin/env bash
export PATH="${gobin_dir}:${homebrew_dir}:${node_dir}:\${PATH}"
cd "${wt_path}"
printf '\e]2;Agent: ${plan_id}\a'
${claude_bin} "\$(cat ${prompt_file})"
LAUNCHER

    log "Opening Ghostty tab for ${plan_id}"
    open_ghostty_tab "${launcher}" "${plan_id}"
    ok "Agent ${plan_id} launched in Ghostty tab"
}

# Wave definitions
wave1() {
    log "=== Wave 1: No dependencies (4 parallel agents) ==="
    launch_agent "01-runner-interface" "devctrl/01-runner-interface" "1"
    launch_agent "02-worktree-identity" "devctrl/02-worktree-identity" "1"
    launch_agent "03-robustness-fixes" "devctrl/03-robustness-fixes" "1"
    launch_agent "04-worktree-exclusions" "devctrl/04-worktree-exclusions" "1"
    log ""
    ok "Wave 1: 4 agents launched in separate Ghostty tabs."
    warn "After all Wave 1 PRs are merged to main, run: $0 wave2"
}

wave2() {
    log "=== Wave 2: Depends on Plans 01, 02 (3 parallel agents) ==="
    warn "Prerequisite: Plans 01 and 02 must be merged to main."
    echo ""
    read -p "Have plans 01 and 02 been merged? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        err "Aborting. Merge Wave 1 PRs first."
        return 1
    fi

    launch_agent "05-error-handling" "devctrl/05-error-handling" "2"
    launch_agent "06-env-provider" "devctrl/06-env-provider" "2"
    launch_agent "08-new-commands" "devctrl/08-new-commands" "2"
    log ""
    ok "Wave 2: 3 agents launched in separate Ghostty tabs."
    warn "After all Wave 2 PRs are merged to main, run: $0 wave3"
}

wave3() {
    log "=== Wave 3: Depends on Plans 01, 02, 06 (1 agent) ==="
    warn "Prerequisite: Plans 01, 02, and 06 must be merged to main."
    echo ""
    read -p "Have plans 01, 02, and 06 been merged? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        err "Aborting. Merge Wave 2 PRs first."
        return 1
    fi

    launch_agent "07-docker-isolation" "devctrl/07-docker-isolation" "3"
    log ""
    ok "Wave 3: 1 agent launched in Ghostty tab."
    warn "After Wave 3 PR is merged to main, run: $0 wave4"
}

wave4() {
    log "=== Wave 4: Depends on Plans 01, 02, 05 (1 agent) ==="
    warn "Prerequisite: Plans 01, 02, and 05 must be merged to main."
    echo ""
    read -p "Have plans 01, 02, and 05 been merged? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        err "Aborting. Merge Wave 3 PRs first."
        return 1
    fi

    launch_agent "09-tests-and-claudemd" "devctrl/09-tests-and-claudemd" "4"
    log ""
    ok "Final wave launched in Ghostty tab. After this PR is merged, the refactor is complete!"
}

show_status() {
    log "=== Worktree Status ==="
    echo ""

    local plans=(
        "01-runner-interface"
        "02-worktree-identity"
        "03-robustness-fixes"
        "04-worktree-exclusions"
        "05-error-handling"
        "06-env-provider"
        "07-docker-isolation"
        "08-new-commands"
        "09-tests-and-claudemd"
    )

    printf "%-30s %-20s %-15s\n" "PLAN" "BRANCH" "WORKTREE"
    printf "%-30s %-20s %-15s\n" "----" "------" "--------"

    for plan in "${plans[@]}"; do
        local branch="devctrl/${plan}"
        local wt="${WORKTREE_DIR}/plan-${plan}"

        local wt_status="missing"
        if [[ -d "${wt}" ]]; then
            wt_status="exists"
        fi

        local branch_status
        if git rev-parse --verify "${branch}" &>/dev/null; then
            local commits
            commits=$(git rev-list --count "main..${branch}" 2>/dev/null || echo "?")
            branch_status="+${commits} commits"
        else
            branch_status="not created"
        fi

        printf "%-30s %-20s %-15s\n" "${plan}" "${branch_status}" "${wt_status}"
    done
}

do_cleanup() {
    log "=== Cleaning up plan worktrees ==="
    warn "This will remove all .worktrees/plan-* directories."
    echo ""
    read -p "Continue? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        return 0
    fi

    local plans=(
        "01-runner-interface"
        "02-worktree-identity"
        "03-robustness-fixes"
        "04-worktree-exclusions"
        "05-error-handling"
        "06-env-provider"
        "07-docker-isolation"
        "08-new-commands"
        "09-tests-and-claudemd"
    )

    for plan in "${plans[@]}"; do
        local wt="${WORKTREE_DIR}/plan-${plan}"
        if [[ -d "${wt}" ]]; then
            log "Removing worktree: ${wt}"
            git worktree remove "${wt}" --force 2>/dev/null || {
                warn "Could not remove ${wt} — may need manual cleanup"
            }
        fi
    done

    # Clean up temp prompt/launcher files
    if [[ -d "${PROMPT_DIR}" ]]; then
        log "Removing temp prompt files: ${PROMPT_DIR}"
        rm -rf "${PROMPT_DIR}"
    fi

    ok "Cleanup complete."
}

# Main
case "${1:-help}" in
    wave1|all)
        wave1
        ;;
    wave2)
        wave2
        ;;
    wave3)
        wave3
        ;;
    wave4)
        wave4
        ;;
    status)
        show_status
        ;;
    cleanup)
        do_cleanup
        ;;
    help|*)
        echo "Usage: $0 {wave1|wave2|wave3|wave4|all|status|cleanup}"
        echo ""
        echo "  wave1    Launch Wave 1 (01-04, no deps, 4 parallel agents)"
        echo "  wave2    Launch Wave 2 (05-06,08, needs Wave 1 merged, 3 agents)"
        echo "  wave3    Launch Wave 3 (07, needs Wave 2 merged, 1 agent)"
        echo "  wave4    Launch Wave 4 (09, needs Wave 3 merged, 1 agent)"
        echo "  all      Alias for wave1"
        echo "  status   Show status of all worktrees and branches"
        echo "  cleanup  Remove all plan worktrees"
        ;;
esac

# No wait needed — agents run interactively in their own Ghostty tabs
