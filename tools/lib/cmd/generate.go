package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

func init() {
	generateMockCmd.AddCommand(generateDBCMockCmd)
	generateMockCmd.AddCommand(generateDAOMockCmd)
	generateMockCmd.AddCommand(generateSMSMockCmd)
	generateCmd.AddCommand(generateProtoCmd)
	generateCmd.AddCommand(generateSQLcCmd)
	generateCmd.AddCommand(generateEnumCmd)
	generateCmd.AddCommand(generateMockCmd)
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().Bool("watch", false, "Watch the input directory and automatically re-run the generator.")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Run all code generation sub-tasks",
	Run: func(cmd *cobra.Command, args []string) {
		runPar(cmd, args,
			generateProtoCmd,
			generateSQLcCmd,
			generateEnumCmd,
		)
		generateMockCmd.Run(cmd, args)
	},
}

var generateProtoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Generate protobuf and gRPC code",
	Run: func(cmd *cobra.Command, _ []string) {
		watchFlag, err := cmd.Flags().GetBool("watch")
		guard(err, "check '--watch' flag")

		guard(generateProto(), "execute `buf ...` command")
		if !watchFlag {
			return
		}

		go watch(filepath.FromSlash("proto/"), "Proto codegen", func(_ watcher.Event) {
			if err := generateProto(); err != nil {
				log.Printf("Encountered error while running 'Proto codegen' task. Waiting to re-run...")
			}
		})

		<-shuttingDown
	},
}

func generateProto() error {
	buf := exec.Command("buf", "generate", "--template", filepath.FromSlash("buf.gen.dev.yaml")) //nolint:gosec // Input is not dynamically provided by end-user.
	buf.Stdout = os.Stdout
	buf.Stderr = os.Stderr

	err := buf.Run()
	if err != nil {
		return fmt.Errorf("run buf generate cmd: %w", err)
	}

	return nil
}

var generateSQLcCmd = &cobra.Command{
	Use:   "sqlc",
	Short: "Generate SQLc queries and data holders",
	Run: func(cmd *cobra.Command, _ []string) {
		watchFlag, err := cmd.Flags().GetBool("watch")
		guard(err, "check '--watch' flag")

		guard(generateSQLc(), "execute `sqlc ...` command")
		if !watchFlag {
			return
		}

		go watch(filepath.FromSlash("api/internal/db/"), "SQLc codegen", func(_ watcher.Event) {
			if err := generateSQLc(); err != nil {
				log.Printf("Encountered error while running 'SQLc codegen' task. Waiting to re-run...")
			}
		})

		<-shuttingDown
	},
}

func generateSQLc() error {
	sqlc := exec.Command("sqlc", "generate", "--file", filepath.FromSlash("api/internal/db/sqlc.yaml")) //nolint:gosec // Input is not dynamically provided by end-user.
	sqlc.Stdout = os.Stdout
	sqlc.Stderr = os.Stderr

	err := sqlc.Run()
	if err != nil {
		return fmt.Errorf("run sqlc generate cmd: %w", err)
	}

	return nil
}

var generateEnumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Generate enum stringers",
	Run: func(_ *cobra.Command, _ []string) {
		guard(generateEnum("ScoringCategory", filepath.FromSlash("./api/internal/lib/models")), "execute `enumer ...` command for ")
	},
}

func generateEnum(typ, pkg string) error {
	enumer := exec.Command("enumer", "-sql", "-transform", "snake-upper", "-type", typ, pkg)
	enumer.Stdout = os.Stdout
	enumer.Stderr = os.Stderr

	err := enumer.Run()
	if err != nil {
		return fmt.Errorf("run enumer cmd: %w", err)
	}

	return nil
}

var generateMockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Generate IO clients' (DAO and SMS) interfaces and mocks",
	Run: func(cmd *cobra.Command, args []string) {
		runPar(cmd, args,
			generateDBCMockCmd,
			generateSMSMockCmd,
		)
		// generateDAOMockCmd must be called after generateDBCMockCmd.
		generateDAOMockCmd.Run(cmd, args)
	},
}

var generateDBCMockCmd = &cobra.Command{
	Use:   "dbc",
	Short: "Generate DBC mock",
	Run: func(_ *cobra.Command, _ []string) {
		guard(
			generateMock("Querier", filepath.FromSlash("api/internal/lib/dao/internal/dbc/")),
			"generate mock DBC",
		)
	},
}

var generateDAOMockCmd = &cobra.Command{
	Use:   "dao",
	Short: "Generate DAO interface and mock",
	Run: func(_ *cobra.Command, _ []string) {
		guard(
			generateInterfaceAndMock(mockConfig{
				targetStruct: "Queries",
				ifaceName:    "QueryProvider",
				ifaceComment: "QueryProvider describes all of the queries exposed by the DAO, to allow for testing mocks.",
				genDir:       filepath.FromSlash("api/internal/lib/dao/"),
				pkgName:      "dao",
			}),
			"generate mock DAO",
		)
	},
}

var generateSMSMockCmd = &cobra.Command{
	Use:   "sms",
	Short: "Generate SMS client interface and mock",
	Run: func(_ *cobra.Command, _ []string) {
		guard(
			generateInterfaceAndMock(mockConfig{
				targetStruct: "Client",
				ifaceName:    "Messenger",
				ifaceComment: "Messenger describes all of the operations exposed by the SMS client, to allow for testing mocks.",
				genDir:       filepath.FromSlash("api/internal/lib/sms/"),
				pkgName:      "sms",
			}),
			"generate mock SMS Client",
		)
	},
}

type mockConfig struct {
	targetStruct string
	ifaceName    string
	ifaceComment string
	genDir       string
	pkgName      string
}

func generateInterfaceAndMock(mc mockConfig) error {
	if err := generateInterface(mc); err != nil {
		return fmt.Errorf("generate interface: %w", err)
	}

	if err := generateMock(mc.ifaceName, mc.genDir); err != nil {
		return fmt.Errorf("generate mock: %w", err)
	}

	return nil
}

func generateInterface(mc mockConfig) error {
	args := []string{
		"--struct", mc.targetStruct,
		"--iface", mc.ifaceName,
		"--iface-comment", mc.ifaceComment,
		"--comment", "// Code generated by ifacemaker; DO NOT EDIT.",
		"--pkg", mc.pkgName,
		"--output", filepath.Join(mc.genDir, "gen_interface.go"),
	}

	files, err := os.ReadDir(mc.genDir)
	if err != nil {
		return fmt.Errorf("read DAO dir: %w", err)
	}

	for _, f := range files {
		// Skip directories, symlinks, etc.
		if !f.Type().IsRegular() {
			continue
		}

		// Skip generated files.
		if strings.HasPrefix(f.Name(), "gen_") {
			continue
		}

		// Skip test files to avoid pulling in anonymous imports required for connecting to the DB.
		if strings.HasSuffix(f.Name(), "_test.go") {
			continue
		}

		args = append(args, "--file", filepath.Join(mc.genDir, f.Name()))
	}

	iface := exec.Command("ifacemaker", args...)

	iface.Stdout = os.Stdout
	iface.Stderr = os.Stderr

	if err := iface.Run(); err != nil {
		return fmt.Errorf("run interface generator command: %w", err)
	}

	return nil
}

func generateMock(ifaceName, genDir string) error {
	mock := exec.Command("mockery",
		"--dir", genDir,
		"--name", ifaceName,
		"--filename", "gen_mock.go",
		"--inpackage",
	)

	mock.Stdout = os.Stdout
	mock.Stderr = os.Stderr

	if err := mock.Run(); err != nil {
		return fmt.Errorf("run mock generator command: %w", err)
	}

	return nil
}
