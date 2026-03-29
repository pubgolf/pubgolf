package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// errMissingTools is returned when doctor finds tools that are not installed.
var errMissingTools = errors.New("missing tools")

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check that all development tools are installed and configured",
	Run: func(cmd *cobra.Command, _ []string) {
		classifyAndExit(runDoctorChecks(cmd.Context()))
	},
}

type doctorCheck struct {
	name    string
	command string
	args    []string
	parse   func(output string) string // extract version from output; nil means use first word
	hint    string                     // help text on failure
}

var doctorChecks = []doctorCheck{
	{
		name:    "Go",
		command: "go",
		args:    []string{"version"},
		parse: func(output string) string {
			// "go version go1.26.0 darwin/arm64" -> "1.26.0"
			for field := range strings.FieldsSeq(output) {
				if after, ok := strings.CutPrefix(field, "go1"); ok {
					return "1" + after
				}
			}

			return firstWord(output)
		},
	},
	{
		name:    "Docker",
		command: "docker",
		args:    []string{"info", "--format", "{{.ServerVersion}}"},
		parse: func(output string) string {
			v := strings.TrimSpace(output)
			if v == "" {
				return ""
			}

			return "running (" + v + ")"
		},
		hint: "run 'open -a Docker' or start dockerd",
	},
	{
		name:    "Doppler",
		command: "doppler",
		args:    []string{"--version"},
		parse: func(output string) string {
			// "v3.68.0" or "doppler v3.68.0"
			for field := range strings.FieldsSeq(output) {
				if strings.HasPrefix(field, "v") {
					return field
				}
			}

			return firstWord(output)
		},
	},
	{
		name:    "golangci-lint",
		command: "golangci-lint",
		args:    []string{"--version"},
		parse: func(output string) string {
			// "golangci-lint has version 2.1.0 built ..."
			_, after, found := strings.Cut(output, "version ")
			if found {
				return firstWord(after)
			}

			return firstWord(output)
		},
	},
	{
		name:    "buf",
		command: "buf",
		args:    []string{"--version"},
	},
	{
		name:    "sqlc",
		command: "sqlc",
		args:    []string{"version"},
	},
	{
		name:    "mockery",
		command: "mockery",
		args:    []string{"--version"},
	},
	{
		name:    "Node.js",
		command: "node",
		args:    []string{"--version"},
	},
	{
		name:    "npm",
		command: "npm",
		args:    []string{"--version"},
	},
}

func runDoctorChecks(ctx context.Context) error {
	w := os.Stdout

	fmt.Fprintln(w, "Environment check:")

	var missing []string

	for _, check := range doctorChecks {
		version, err := runDoctorCheck(ctx, check)
		if err != nil {
			missing = append(missing, check.name)

			hint := ""
			if check.hint != "" {
				hint = "  (" + check.hint + ")"
			}

			fmt.Fprintf(w, "  %-16s %-18s %s%s\n", check.name+":", "not found", "\u2717", hint)
		} else {
			fmt.Fprintf(w, "  %-16s %-18s %s\n", check.name+":", version, "\u2713")
		}
	}

	// Project root check.
	fmt.Fprintf(w, "  %-16s %-18s %s\n", "Project root:", projectRoot, "\u2713")

	// Web-app npm dependencies check.
	nodeModulesPath := filepath.Join(projectRoot, "web-app", "node_modules")

	_, statErr := os.Stat(nodeModulesPath)
	if statErr != nil {
		missing = append(missing, "web-app npm deps")

		fmt.Fprintf(w, "  %-16s %-18s %s  (run 'pubgolf-devctrl install web')\n", "npm deps:", "not installed", "\u2717")
	} else {
		fmt.Fprintf(w, "  %-16s %-18s %s\n", "npm deps:", "installed", "\u2713")
	}

	if len(missing) > 0 {
		return fmt.Errorf("%w: %s", errMissingTools, strings.Join(missing, ", "))
	}

	return nil
}

func runDoctorCheck(ctx context.Context, check doctorCheck) (string, error) {
	out, err := exec.CommandContext(ctx, check.command, check.args...).Output() //nolint:gosec // Commands are hardcoded check definitions.
	if err != nil {
		return "", fmt.Errorf("run %s: %w", check.command, err)
	}

	output := strings.TrimSpace(string(out))

	if check.parse != nil {
		return check.parse(output), nil
	}

	return firstWord(output), nil
}

func firstWord(s string) string {
	s = strings.TrimSpace(s)

	if before, _, found := strings.Cut(s, " "); found {
		return before
	}

	return s
}
