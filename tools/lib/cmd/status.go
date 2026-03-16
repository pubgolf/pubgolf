package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	statusCmd.Flags().Bool("all", false, "Show status for all worktrees")
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show worktree identity, port assignments, and service status",
	Run: func(cmd *cobra.Command, _ []string) {
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag {
			classifyAndExit(printAllWorktreeStatus(cmd.Context()))
		} else {
			classifyAndExit(printCurrentStatus(cmd.Context()))
		}
	},
}

func printCurrentStatus(ctx context.Context) error {
	slug, err := worktreeSlug(ctx)
	if err != nil {
		return fmtErr(err, "determine worktree slug")
	}

	offset, err := portOffsetForSlug(slug)
	if err != nil {
		return fmtErr(err, "compute port offset")
	}

	project := dockerProjectForSlug(slug)

	displaySlug := slug
	if displaySlug == "" {
		displaySlug = "(main working tree)"
	}

	w := os.Stdout

	fmt.Fprintf(w, "Worktree:        %s\n", displaySlug)
	fmt.Fprintf(w, "Port offset:     %d\n", offset)
	fmt.Fprintf(w, "Docker project:  %s\n", project)
	fmt.Fprintf(w, "DB port:         %d\n", 5432+offset)
	fmt.Fprintf(w, "API port:        %d\n", 5000+offset)
	fmt.Fprintf(w, "Minio port:      %d\n", 9000+offset)
	fmt.Fprintf(w, "DB volume:       %s\n", dataDirForSlug("./data/postgres", slug))
	fmt.Fprintf(w, "Minio volume:    %s\n", dataDirForSlug("./data/minio", slug))
	fmt.Fprintf(w, "Coverage dir:    %s\n", dataDirForSlug("./data/go-test-coverage", slug))

	status := queryDockerStatus(ctx, project)
	fmt.Fprintf(w, "Docker status:   %s\n", status)

	return nil
}

func printAllWorktreeStatus(ctx context.Context) error {
	worktrees, err := listWorktrees(ctx)
	if err != nil {
		return fmtErr(err, "list worktrees")
	}

	w := os.Stdout

	fmt.Fprintf(w, "%-20s %-22s %7s %8s %9s  %s\n", "WORKTREE", "SLUG", "OFFSET", "DB_PORT", "API_PORT", "DOCKER")

	for _, wt := range worktrees {
		offset, _ := portOffsetForSlug(wt.slug)
		project := dockerProjectForSlug(wt.slug)
		status := queryDockerStatus(ctx, project)

		displayName := wt.name
		displaySlug := wt.slug

		if wt.slug == "" {
			displayName = "(main)"
			displaySlug = "-"
		}

		fmt.Fprintf(w, "%-20s %-22s %7d %8d %9d  %s\n",
			truncate(displayName, 20),
			truncate(displaySlug, 22),
			offset,
			5432+offset,
			5000+offset,
			status,
		)
	}

	return nil
}

type worktreeInfo struct {
	name   string // directory basename
	slug   string // normalized slug ("" for main)
	path   string // absolute path to worktree
	branch string // branch name (e.g. "devctrl/08-new-commands"); empty for detached HEAD
}

// listWorktrees parses `git worktree list --porcelain` to enumerate all worktrees.
func listWorktrees(ctx context.Context) ([]worktreeInfo, error) {
	out, err := exec.CommandContext(ctx, "git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}

	// Find the main worktree by resolving the common git dir.
	commonDirOut, err := exec.CommandContext(ctx, "git", "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return nil, fmt.Errorf("git rev-parse --git-common-dir: %w", err)
	}

	commonDir, err := filepath.Abs(strings.TrimSpace(string(commonDirOut)))
	if err != nil {
		return nil, fmt.Errorf("resolve common dir: %w", err)
	}

	// The main worktree is the parent of the common .git dir.
	mainWorktreePath := filepath.Dir(commonDir)

	var worktrees []worktreeInfo

	for block := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n\n") {
		var wtPath, branch string

		for line := range strings.SplitSeq(block, "\n") {
			if after, ok := strings.CutPrefix(line, "worktree "); ok {
				wtPath = after
			}

			if after, ok := strings.CutPrefix(line, "branch refs/heads/"); ok {
				branch = after
			}
		}

		if wtPath == "" {
			continue
		}

		absPath, _ := filepath.Abs(wtPath)

		if absPath == mainWorktreePath {
			worktrees = append(worktrees, worktreeInfo{
				name:   "(main)",
				slug:   "",
				path:   absPath,
				branch: branch,
			})
		} else {
			name := filepath.Base(absPath)
			worktrees = append(worktrees, worktreeInfo{
				name:   name,
				slug:   normalizeSlug(name),
				path:   absPath,
				branch: branch,
			})
		}
	}

	return worktrees, nil
}

// queryDockerStatus checks Docker container status for a given project name.
func queryDockerStatus(ctx context.Context, project string) string {
	out, err := exec.CommandContext(ctx, "docker", "compose", "ls", //nolint:gosec // Project name is derived from worktree slug, not user input.
		"--all",
		"--format", "json",
		"--filter", "name="+project,
	).Output()
	if err != nil {
		return "unknown (Docker unavailable)"
	}

	var projects []struct {
		Name   string `json:"Name"`
		Status string `json:"Status"`
	}

	unmarshalErr := json.Unmarshal(bytes.TrimSpace(out), &projects)
	if unmarshalErr != nil {
		return "unknown"
	}

	for _, p := range projects {
		if p.Name == project {
			return p.Status
		}
	}

	return "not running"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-2] + ".."
}
