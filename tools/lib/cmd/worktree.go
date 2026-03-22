package cmd

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// errEmptySlug is returned when a worktree directory name normalizes to an empty string.
var errEmptySlug = errors.New("worktree directory normalizes to an empty slug")

// worktreeSlug returns a normalized identifier for the current git worktree.
// Returns ("", nil) for the main working tree, (slug, nil) for a worktree,
// or ("", err) if git commands fail.
func worktreeSlug() (string, error) {
	ctx := context.Background()

	topLevelOut, err := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}

	commonDirOut, err := exec.CommandContext(ctx, "git", "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --git-common-dir: %w", err)
	}

	topLevel, err := filepath.Abs(strings.TrimSpace(string(topLevelOut)))
	if err != nil {
		return "", fmt.Errorf("resolve toplevel path: %w", err)
	}

	commonDir, err := filepath.Abs(strings.TrimSpace(string(commonDirOut)))
	if err != nil {
		return "", fmt.Errorf("resolve common dir path: %w", err)
	}

	// If commonDir is not a prefix of topLevel, we're in a worktree.
	if !strings.HasPrefix(commonDir, topLevel) {
		raw := filepath.Base(topLevel)
		slug := normalizeSlug(raw)

		if slug == "" {
			return "", fmt.Errorf("%w: %q", errEmptySlug, raw)
		}

		return slug, nil
	}

	return "", nil
}

// nonAlphanumeric matches any character that is not a lowercase letter or digit.
var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

// normalizeSlug converts a raw worktree directory name into a clean slug.
// It lowercases the input, replaces non-alphanumeric runs with single hyphens,
// trims leading/trailing hyphens, and truncates long names with a hash suffix.
func normalizeSlug(raw string) string {
	s := strings.ToLower(raw)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	if len(s) > 20 {
		h := fnv.New32a()
		h.Write([]byte(s))
		hash := fmt.Sprintf("%06x", h.Sum32())[:6]
		s = s[:20] + "-" + hash
	}

	return s
}

// portOffsetForSlug computes the port offset for a given slug.
// Returns 0 for empty slug (main tree). For non-empty slugs, checks
// PUBGOLF_PORT_OFFSET env var first, then falls back to FNV-32a hash.
// The offset is always in the range [1, 500].
func portOffsetForSlug(slug string) int {
	if slug == "" {
		return 0
	}

	if v := os.Getenv("PUBGOLF_PORT_OFFSET"); v != "" {
		offset, err := strconv.Atoi(v)
		if err == nil && offset > 0 && offset < 500 {
			return offset
		}
	}

	h := fnv.New32a()
	h.Write([]byte(slug))

	return int(h.Sum32()%500) + 1
}

// worktreePortOffset returns the port offset for the current worktree.
func worktreePortOffset() int {
	slug, _ := worktreeSlug()

	return portOffsetForSlug(slug)
}

// dockerProjectForSlug returns the Docker Compose project name for a given slug.
// Returns "pubgolf" for empty slug (main tree), "pubgolf-<slug>" for worktrees.
func dockerProjectForSlug(slug string) string {
	if slug == "" {
		return "pubgolf"
	}

	return "pubgolf-" + slug
}

// worktreeDockerProject returns the Docker Compose project name for the current worktree.
func worktreeDockerProject() string {
	slug, _ := worktreeSlug()

	return dockerProjectForSlug(slug)
}

// dataDirForSlug returns the data directory path for a given slug.
// Returns base for empty slug (main tree), "base-<slug>" for worktrees.
func dataDirForSlug(base, slug string) string {
	if slug == "" {
		return base
	}

	return base + "-" + slug
}

// worktreeDataDir returns the data directory path, namespaced by worktree slug.
func worktreeDataDir(base string) string {
	slug, _ := worktreeSlug()

	return dataDirForSlug(base, slug)
}
