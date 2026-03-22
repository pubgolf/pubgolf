package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Cmd describes a subprocess invocation.
type Cmd struct {
	Name   string
	Args   []string
	Dir    string    // working directory; empty uses process default
	Env    []string  // additional env vars (KEY=VALUE); nil inherits parent env
	Stdin  io.Reader // nil = os.Stdin
	Stdout io.Writer // nil = os.Stdout
	Stderr io.Writer // nil = os.Stderr
}

// String returns a shell-quoted representation of the command, suitable for
// copy-pasting into a terminal.
func (c Cmd) String() string {
	parts := make([]string, 0, 1+len(c.Args))
	parts = append(parts, c.Name)

	for _, arg := range c.Args {
		if needsQuoting(arg) {
			parts = append(parts, "'"+strings.ReplaceAll(arg, "'", "'\\''")+"'")
		} else {
			parts = append(parts, arg)
		}
	}

	return strings.Join(parts, " ")
}

func needsQuoting(s string) bool {
	if s == "" {
		return true
	}

	for _, c := range s {
		if !isShellSafe(c) {
			return true
		}
	}

	return false
}

func isShellSafe(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_' || c == '.' || c == '/' || c == ':' || c == '=' || c == ','
}

// Process represents a running subprocess started via Runner.Start.
type Process interface {
	// Wait blocks until the process exits and returns its error.
	Wait() error
	// Stop sends SIGINT to the process group for graceful shutdown.
	Stop()
}

// Runner executes or records subprocess invocations.
type Runner interface {
	// Run executes (or records) a single command and waits for completion.
	Run(ctx context.Context, cmd Cmd) error
	// Start begins a long-running process and returns a Process handle.
	Start(ctx context.Context, cmd Cmd) (Process, error)
}

// ExecRunner is the production Runner implementation that executes real subprocesses.
type ExecRunner struct{}

// Run executes a command and waits for completion.
func (r ExecRunner) Run(ctx context.Context, c Cmd) error {
	logCmd("exec", c)

	cmd := r.buildCmd(ctx, c)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("run %s: %w", c.Name, err)
	}

	return nil
}

// Start begins a long-running process with its own process group.
func (r ExecRunner) Start(ctx context.Context, c Cmd) (Process, error) { //nolint:ireturn // Interface return is intentional for Runner contract.
	logCmd("exec", c)

	cmd := r.buildCmd(ctx, c)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start %s: %w", c.Name, err)
	}

	return &execProcess{cmd: cmd}, nil
}

func (r ExecRunner) buildCmd(ctx context.Context, c Cmd) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.Name, c.Args...) //nolint:gosec // Command construction is controlled by devctrl internals.

	if c.Dir != "" {
		cmd.Dir = c.Dir
	}

	if c.Stdin != nil {
		cmd.Stdin = c.Stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	if c.Stdout != nil {
		cmd.Stdout = c.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}

	if c.Stderr != nil {
		cmd.Stderr = c.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	if len(c.Env) > 0 {
		cmd.Env = mergeEnv(os.Environ(), c.Env)
	}

	return cmd
}

type execProcess struct {
	cmd *exec.Cmd
}

func (p *execProcess) Wait() error {
	err := p.cmd.Wait()
	if err != nil {
		return fmt.Errorf("wait %s: %w", p.cmd.Path, err)
	}

	return nil
}

func (p *execProcess) Stop() {
	if p.cmd.Process == nil {
		return
	}

	pgid, err := syscall.Getpgid(p.cmd.Process.Pid)
	if err != nil {
		return
	}

	_ = syscall.Kill(-pgid, syscall.SIGINT)
}

// DryRunner records commands without executing them. Used for --dry-run and tests.
type DryRunner struct {
	Recorded []Cmd
	ErrorFor map[string]error // keyed on Cmd.Name
}

// Run records the command and returns the configured error (or nil).
func (r *DryRunner) Run(_ context.Context, c Cmd) error {
	r.Recorded = append(r.Recorded, c)
	logCmd("dry-run", c)

	if r.ErrorFor != nil {
		if err, ok := r.ErrorFor[c.Name]; ok {
			return err
		}
	}

	return nil
}

// Start records the command and returns a no-op Process.
func (r *DryRunner) Start(_ context.Context, c Cmd) (Process, error) { //nolint:ireturn // Interface return is intentional for Runner contract.
	r.Recorded = append(r.Recorded, c)
	logCmd("dry-run", c)

	return &dryProcess{}, nil
}

type dryProcess struct{}

func (p *dryProcess) Wait() error { return nil }
func (p *dryProcess) Stop()       {}

// logCmd logs a command invocation with the given prefix.
func logCmd(prefix string, c Cmd) {
	log.Printf("%s: %s", prefix, c.String())

	if c.Dir != "" {
		log.Printf("  dir: %s", c.Dir)
	}

	if len(c.Env) > 0 {
		log.Printf("  env: %s", strings.Join(c.Env, " "))
	}
}

// mergeEnv merges additional env vars into the parent environment.
// Later values override earlier ones with the same key.
func mergeEnv(parent, extra []string) []string {
	m := make(map[string]string, len(parent)+len(extra))
	order := make([]string, 0, len(parent)+len(extra))

	addToMap := func(vars []string) {
		for _, v := range vars {
			k, _, _ := strings.Cut(v, "=")
			if _, exists := m[k]; !exists {
				order = append(order, k)
			}

			m[k] = v
		}
	}

	addToMap(parent)
	addToMap(extra)

	result := make([]string, 0, len(order))
	for _, k := range order {
		result = append(result, m[k])
	}

	return result
}

// fmtErr wraps an error with a descriptive message, matching the existing guard pattern.
func fmtErr(err error, msg string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}
