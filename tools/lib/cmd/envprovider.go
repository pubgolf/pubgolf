package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// EnvProvider resolves environment variables for subprocess execution.
type EnvProvider interface {
	// Env returns environment variables for a given project and config.
	// The returned slice contains KEY=VALUE pairs.
	Env(ctx context.Context, project, config string) ([]string, error)
}

// DopplerProvider fetches environment variables from the Doppler CLI.
type DopplerProvider struct {
	Runner Runner
}

// Env calls `doppler secrets --project P --config C --json` and parses the
// JSON output into KEY=VALUE pairs.
func (p *DopplerProvider) Env(ctx context.Context, project, cfg string) ([]string, error) {
	var buf bytes.Buffer

	err := p.Runner.Run(ctx, Cmd{
		Name:   "doppler",
		Args:   []string{"secrets", "--project", project, "--config", cfg, "--json"},
		Stdout: &buf,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch doppler secrets: %w", err)
	}

	var data map[string]any

	decodeErr := json.NewDecoder(&buf).Decode(&data)
	if decodeErr != nil {
		return nil, fmt.Errorf("parse doppler JSON: %w", decodeErr)
	}

	env := make([]string, 0, len(data))

	for key, secret := range data {
		inner, ok := secret.(map[string]any)
		if !ok {
			continue
		}

		val, ok := inner["computed"]
		if !ok {
			continue
		}

		v, ok := val.(string)
		if !ok {
			continue
		}

		env = append(env, key+"="+v)
	}

	sort.Strings(env)

	return env, nil
}

// EnvFileProvider reads environment variables from .env files,
// falling back to the current process environment.
type EnvFileProvider struct {
	// ProjectRoot is the directory to look for .env files in.
	ProjectRoot string
}

// Env reads .env.{config} from the project root. For keys not found in the
// file, the current process environment is used as a fallback.
func (p *EnvFileProvider) Env(_ context.Context, _, cfg string) ([]string, error) {
	envFile := filepath.Join(p.ProjectRoot, ".env."+cfg)

	fileVars, err := parseEnvFile(envFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("read env file %s: %w", envFile, err)
	}

	// Start with process environment, then overlay file values.
	merged := mergeEnv(os.Environ(), fileVars)

	return merged, nil
}

// parseEnvFile reads a simple KEY=VALUE env file. It supports blank lines,
// comment lines (starting with #), and optional quoting of values.
func parseEnvFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	var env []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = unquote(val)

		env = append(env, key+"="+val)
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, fmt.Errorf("scan env file: %w", scanErr)
	}

	return env, nil
}

// unquote removes matching single or double quotes from the value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}

	return s
}

// AutoProvider selects DopplerProvider if the `doppler` binary is in PATH,
// otherwise falls back to EnvFileProvider.
type AutoProvider struct {
	Runner      Runner
	ProjectRoot string
}

// Env checks for doppler availability and delegates accordingly.
func (p *AutoProvider) Env(ctx context.Context, project, cfg string) ([]string, error) {
	_, pathErr := lookPath("doppler")
	if pathErr == nil {
		dp := &DopplerProvider{Runner: p.Runner}

		return dp.Env(ctx, project, cfg)
	}

	log.Println("Doppler not available, using environment variables.")

	ep := &EnvFileProvider{ProjectRoot: p.ProjectRoot}

	return ep.Env(ctx, project, cfg)
}

// lookPath is a variable so tests can override it.
var lookPath = defaultLookPath

func defaultLookPath(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("look up %s in PATH: %w", name, err)
	}

	return path, nil
}

// newEnvProvider creates an EnvProvider based on the flag value.
func newEnvProvider(flag string, dryRun bool, r Runner) EnvProvider { //nolint:ireturn // Interface return is intentional for provider selection.
	if dryRun {
		return &dryRunProvider{}
	}

	switch flag {
	case "doppler":
		return &DopplerProvider{Runner: r}
	case "env":
		return &EnvFileProvider{ProjectRoot: projectRoot}
	default: // "auto"
		return &AutoProvider{Runner: r, ProjectRoot: projectRoot}
	}
}

// dryRunProvider records the intent to fetch secrets without executing anything.
// In dry-run mode, it returns an empty slice so downstream commands use defaults.
type dryRunProvider struct{}

func (p *dryRunProvider) Env(_ context.Context, project, cfg string) ([]string, error) {
	log.Printf("dry-run: would fetch env for project=%s config=%s", project, cfg)

	return nil, nil
}
