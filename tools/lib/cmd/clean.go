package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().Bool("force", false, "Actually remove orphaned resources (default is dry-run)")
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove merged worktrees and orphaned resources (Docker containers, data dirs)",
	Run: func(cmd *cobra.Command, _ []string) {
		force, _ := cmd.Flags().GetBool("force")
		classifyAndExit(cleanOrphans(cmd.Context(), force))
	},
}

func cleanOrphans(ctx context.Context, force bool) error {
	if !force {
		log.Println("Dry run — pass --force to remove these resources:")
	}

	found := false

	// Pass 1: Remove worktrees whose branches have been merged into main.
	// This runs first so that newly-orphaned Docker/filesystem resources
	// are caught by the subsequent passes.
	found = cleanMergedWorktrees(ctx, force) || found

	// Recompute active slugs after worktree removal.
	activeSlugs, err := activeWorktreeSlugs(ctx)
	if err != nil {
		return fmtErr(err, "enumerate worktrees")
	}

	// Pass 2: Blob storage bucket orphans (must run before Docker pass,
	// which may stop the shared Minio instance).
	found = cleanBucketOrphans(ctx, activeSlugs, force) || found

	// Pass 3: Docker orphans.
	found = cleanDockerOrphans(ctx, activeSlugs, force) || found

	// Pass 4: Filesystem orphans.
	found = cleanFilesystemOrphans(activeSlugs, force) || found

	if !found {
		log.Println("No orphaned resources found.")
	}

	return nil
}

// activeWorktreeSlugs returns the set of slugs for all active worktrees.
// The main worktree has slug "".
func activeWorktreeSlugs(ctx context.Context) (map[string]bool, error) {
	worktrees, err := listWorktrees(ctx)
	if err != nil {
		return nil, err
	}

	slugs := make(map[string]bool, len(worktrees))
	for _, wt := range worktrees {
		slugs[wt.slug] = true
	}

	return slugs, nil
}

// mergedBranches returns the set of branch names that have been merged into main.
func mergedBranches(ctx context.Context) (map[string]bool, error) {
	out, err := exec.CommandContext(ctx, "git", "branch", "--merged", "main").Output()
	if err != nil {
		return nil, fmt.Errorf("git branch --merged main: %w", err)
	}

	branches := make(map[string]bool)

	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		name := strings.TrimSpace(line)

		// Skip the current branch marker and main itself.
		if strings.HasPrefix(name, "* ") || name == "main" || name == "" {
			continue
		}

		branches[name] = true
	}

	return branches, nil
}

// cleanMergedWorktrees removes worktrees whose branches have been merged into main.
// Returns true if any merged worktrees were found.
func cleanMergedWorktrees(ctx context.Context, force bool) bool {
	merged, err := mergedBranches(ctx)
	if err != nil {
		log.Printf("WARNING: failed to check merged branches: %v", err)

		return false
	}

	if len(merged) == 0 {
		return false
	}

	worktrees, err := listWorktrees(ctx)
	if err != nil {
		log.Printf("WARNING: failed to list worktrees: %v", err)

		return false
	}

	found := false

	for _, wt := range worktrees {
		// Never touch main worktree or worktrees with no branch (detached HEAD).
		if wt.slug == "" || wt.branch == "" {
			continue
		}

		// Skip the worktree we're currently running from.
		if wt.path == projectRoot {
			continue
		}

		if !merged[wt.branch] {
			continue
		}

		found = true

		if force {
			log.Printf("Removing worktree: %s (branch %s, merged into main)", wt.path, wt.branch)

			rmWt := exec.CommandContext(ctx, "git", "worktree", "remove", wt.path) //nolint:gosec // Path from git worktree list output.
			rmWt.Stdout = os.Stdout
			rmWt.Stderr = os.Stderr

			rmErr := rmWt.Run()
			if rmErr != nil {
				log.Printf("WARNING: failed to remove worktree %s: %v", wt.path, rmErr)

				continue
			}

			rmBranch := exec.CommandContext(ctx, "git", "branch", "-d", wt.branch) //nolint:gosec // Branch name from git worktree list output.
			rmBranch.Stdout = os.Stdout
			rmBranch.Stderr = os.Stderr

			branchErr := rmBranch.Run()
			if branchErr != nil {
				log.Printf("WARNING: failed to delete branch %s: %v", wt.branch, branchErr)
			}
		} else {
			fmt.Fprintf(os.Stdout, "  [merged] worktree: %s (branch %s)\n", wt.path, wt.branch)
		}
	}

	return found
}

// cleanBucketOrphans finds and removes Minio buckets matching pubgolf-dev-*
// that don't correspond to active worktrees. Returns true if any orphans were found.
func cleanBucketOrphans(ctx context.Context, activeSlugs map[string]bool, force bool) bool {
	buckets, err := listDevBuckets(ctx, envProvider)
	if err != nil {
		log.Println("WARNING: Minio is not available. Skipping bucket cleanup.")
		log.Println("  Start Minio and re-run to clean up orphaned buckets.")

		return false
	}

	found := false

	for _, bucket := range buckets {
		slug := slugFromBucket(bucket)

		if activeSlugs[slug] {
			continue // Active worktree.
		}

		// Never remove the main bucket.
		if bucket == blobBucketPrefix {
			continue
		}

		found = true

		if force {
			log.Printf("Removing blob bucket: %s", bucket)

			rmErr := deleteBucket(ctx, envProvider, bucket)
			if rmErr != nil {
				log.Printf("WARNING: failed to remove bucket %s: %v", bucket, rmErr)
			}
		} else {
			fmt.Fprintf(os.Stdout, "  [orphan] blob bucket: %s\n", bucket)
		}
	}

	return found
}

// cleanDockerOrphans finds and removes Docker Compose projects matching pubgolf-*
// that don't correspond to active worktrees. Returns true if any orphans were found.
func cleanDockerOrphans(ctx context.Context, activeSlugs map[string]bool, force bool) bool {
	out, err := exec.CommandContext(ctx, "docker", "compose", "ls",
		"--all",
		"--format", "json",
	).Output()
	if err != nil {
		log.Println("WARNING: Docker daemon is not available. Skipping container cleanup.")
		log.Println("  Start Docker and re-run to clean up Docker resources.")

		return false
	}

	var projects []struct {
		Name   string `json:"Name"`
		Status string `json:"Status"`
	}

	unmarshalErr := json.Unmarshal(bytes.TrimSpace(out), &projects)
	if unmarshalErr != nil {
		log.Printf("WARNING: failed to parse Docker project list: %v", unmarshalErr)

		return false
	}

	found := false

	for _, p := range projects {
		slug, ok := strings.CutPrefix(p.Name, "pubgolf-")
		if !ok {
			continue // Not a pubgolf project.
		}

		if activeSlugs[slug] {
			continue // Active worktree.
		}

		found = true

		if force {
			log.Printf("Removing docker project: %s", p.Name)

			rmCmd := exec.CommandContext(ctx, "docker", "compose", //nolint:gosec // Project name from docker ls output, not user input.
				"--project-name", p.Name,
				"down", "--volumes",
			)
			rmCmd.Stdout = os.Stdout
			rmCmd.Stderr = os.Stderr

			rmErr := rmCmd.Run()
			if rmErr != nil {
				log.Printf("WARNING: failed to remove docker project %s: %v", p.Name, rmErr)
			}
		} else {
			fmt.Fprintf(os.Stdout, "  [orphan] docker project: %s (%s)\n", p.Name, p.Status)
		}
	}

	return found
}

// cleanFilesystemOrphans finds and removes orphaned data directories. Returns true if any orphans were found.
func cleanFilesystemOrphans(activeSlugs map[string]bool, force bool) bool {
	found := false

	for _, base := range []string{"data/postgres", "data/go-test-coverage"} {
		pattern := filepath.Join(projectRoot, base+"-*")

		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("WARNING: failed to glob %s: %v", pattern, err)

			continue
		}

		for _, match := range matches {
			dirName := filepath.Base(match)

			// Extract slug from directory name (e.g., "postgres-fix-auth" -> "fix-auth").
			prefix := filepath.Base(base) + "-"

			slug, ok := strings.CutPrefix(dirName, prefix)
			if !ok {
				continue
			}

			if activeSlugs[slug] {
				continue
			}

			found = true

			size := dirSize(match)

			if force {
				log.Printf("Removing data directory: %s (%s)", match, size)

				rmErr := os.RemoveAll(match)
				if rmErr != nil {
					log.Printf("WARNING: failed to remove %s: %v", match, rmErr)
				}
			} else {
				fmt.Fprintf(os.Stdout, "  [orphan] data directory: %s (%s)\n", dirName, size)
			}
		}
	}

	return found
}

// dirSize returns a human-readable size for a directory.
func dirSize(path string) string {
	var total int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if !info.IsDir() {
			total += info.Size()
		}

		return nil
	})
	if err != nil {
		return "unknown"
	}

	return formatBytes(total)
}

func formatBytes(b int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case b >= gb:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(gb))
	case b >= mb:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(mb))
	case b >= kb:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(kb))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
