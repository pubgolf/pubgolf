package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	stopCmd.Flags().Bool("all", false, "Stop services across all worktree projects")
	stopCmd.Flags().Bool("remove-data", false, "Remove this worktree's data directories after stopping")
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all background processes started with `devctrl run ...`",
	Run: func(cmd *cobra.Command, _ []string) {
		allFlag, _ := cmd.Flags().GetBool("all")
		removeDataFlag, _ := cmd.Flags().GetBool("remove-data")

		// Delete the worktree's blob bucket while Minio is still running.
		if removeDataFlag {
			slug, _ := worktreeSlug(cmd.Context())
			bucket := blobBucketForSlug(slug)

			bucketErr := deleteBucket(cmd.Context(), envProvider, bucket)
			if bucketErr != nil {
				log.Printf("Warning: could not delete blob bucket %q: %v", bucket, bucketErr)
			}
		}

		if allFlag {
			classifyAndExit(stopAllWorktrees(cmd.Context(), runner))
		} else {
			classifyAndExit(dockerStop(cmd.Context(), runner))
		}

		if removeDataFlag {
			classifyAndExit(removeWorktreeData(cmd.Context()))
		}
	},
}

func dockerStop(ctx context.Context, r Runner) error {
	projectName := worktreeDockerProject(ctx)

	// docker compose down doesn't need secrets; run without env injection.
	err := r.Run(ctx, Cmd{
		Name: "docker",
		Args: []string{
			"compose",
			"--file", filepath.FromSlash("./infra/docker-compose.dev.yaml"),
			"--project-name", projectName,
			"down",
		},
	})
	if err != nil {
		return fmtErr(err, "run docker-compose down cmd")
	}

	return nil
}

// stopAllWorktrees stops Docker services for every known worktree project.
func stopAllWorktrees(ctx context.Context, r Runner) error {
	worktrees, err := listWorktrees(ctx)
	if err != nil {
		return fmtErr(err, "list worktrees")
	}

	var stopErrors []error

	for _, wt := range worktrees {
		project := dockerProjectForSlug(wt.slug)
		log.Printf("Stopping project %s...", project)

		stopErr := r.Run(ctx, Cmd{
			Name: "docker",
			Args: []string{
				"compose",
				"--project-name", project,
				"down",
			},
		})
		if stopErr != nil {
			stopErrors = append(stopErrors, fmt.Errorf("stop %s: %w", project, stopErr))
		}
	}

	return errors.Join(stopErrors...)
}

// removeWorktreeData removes this worktree's data directories.
func removeWorktreeData(ctx context.Context) error {
	slug, err := worktreeSlug(ctx)
	if err != nil {
		return fmtErr(err, "determine worktree slug")
	}

	for _, base := range []string{"data/postgres", "data/go-test-coverage"} {
		dir := filepath.Join(projectRoot, dataDirForSlug(base, slug))

		info, statErr := os.Stat(dir)
		if statErr != nil {
			continue // Doesn't exist, nothing to do.
		}

		if !info.IsDir() {
			continue
		}

		log.Printf("Removing %s...", dir)

		rmErr := os.RemoveAll(dir)
		if rmErr != nil {
			return fmtErr(rmErr, "remove "+dir)
		}
	}

	return nil
}
