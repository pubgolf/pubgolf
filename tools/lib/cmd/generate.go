package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
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

		guard(generateProto(cmd.Context(), runner), "execute `buf ...` command")

		if !watchFlag {
			return
		}

		go watch(filepath.FromSlash("proto/"), "Proto codegen", func(_ watcher.Event) {
			err := generateProto(context.Background(), runner)
			if err != nil {
				log.Printf("Encountered error while running 'Proto codegen' task. Waiting to re-run...")
			}
		})

		<-shuttingDown
	},
}

func generateProto(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: "buf",
		Args: []string{"generate", "--template", filepath.FromSlash("buf.gen.dev.yaml")},
	})
	if err != nil {
		return fmtErr(err, "run buf generate cmd")
	}

	return nil
}

var generateSQLcCmd = &cobra.Command{
	Use:   "sqlc",
	Short: "Generate SQLc queries and data holders",
	Run: func(cmd *cobra.Command, _ []string) {
		watchFlag, err := cmd.Flags().GetBool("watch")
		guard(err, "check '--watch' flag")

		guard(generateSQLc(cmd.Context(), runner), "execute `sqlc ...` command")

		if !watchFlag {
			return
		}

		go watch(filepath.FromSlash("api/internal/db/"), "SQLc codegen", func(_ watcher.Event) {
			err := generateSQLc(context.Background(), runner)
			if err != nil {
				log.Printf("Encountered error while running 'SQLc codegen' task. Waiting to re-run...")
			}
		})

		<-shuttingDown
	},
}

func generateSQLc(ctx context.Context, r Runner) error {
	err := r.Run(ctx, Cmd{
		Name: "sqlc",
		Args: []string{"generate", "--file", filepath.FromSlash("api/internal/db/sqlc.yaml")},
	})
	if err != nil {
		return fmtErr(err, "run sqlc generate cmd")
	}

	return nil
}

var generateEnumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Generate enum stringers",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(generateEnum(cmd.Context(), runner, "ScoringCategory", filepath.FromSlash("./api/internal/lib/models")), "execute `enumer ...` command for ")
	},
}

func generateEnum(ctx context.Context, r Runner, typ, pkg string) error {
	err := r.Run(ctx, Cmd{
		Name: "enumer",
		Args: []string{"-sql", "-transform", "snake-upper", "-type", typ, pkg},
	})
	if err != nil {
		return fmtErr(err, "run enumer cmd")
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
	Run: func(cmd *cobra.Command, _ []string) {
		guard(
			generateMock(cmd.Context(), runner, "Querier", filepath.FromSlash("api/internal/lib/dao/internal/dbc/")),
			"generate mock DBC",
		)
	},
}

var generateDAOMockCmd = &cobra.Command{
	Use:   "dao",
	Short: "Generate DAO interface and mock",
	Run: func(cmd *cobra.Command, _ []string) {
		guard(
			generateInterfaceAndMock(cmd.Context(), runner, mockConfig{
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
	Run: func(cmd *cobra.Command, _ []string) {
		guard(
			generateInterfaceAndMock(cmd.Context(), runner, mockConfig{
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

func generateInterfaceAndMock(ctx context.Context, r Runner, mc mockConfig) error {
	err := generateInterface(ctx, r, mc)
	if err != nil {
		return fmt.Errorf("generate interface: %w", err)
	}

	err = generateMock(ctx, r, mc.ifaceName, mc.genDir)
	if err != nil {
		return fmt.Errorf("generate mock: %w", err)
	}

	return nil
}

func generateInterface(ctx context.Context, r Runner, mc mockConfig) error {
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

	err = r.Run(ctx, Cmd{
		Name: "ifacemaker",
		Args: args,
	})
	if err != nil {
		return fmtErr(err, "run interface generator command")
	}

	return nil
}

func generateMock(ctx context.Context, r Runner, ifaceName, genDir string) error {
	err := r.Run(ctx, Cmd{
		Name: "mockery",
		Args: []string{
			"--dir", genDir,
			"--name", ifaceName,
			"--filename", "gen_mock.go",
			"--inpackage",
		},
	})
	if err != nil {
		return fmtErr(err, "run mock generator command")
	}

	return nil
}
