package cmd

import (
	"bytes"
	"database/sql"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	// Postgres and file drivers for the migration lib.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateFixCmd)
	rootCmd.AddCommand(migrateCmd)
}

var (
	migrationDirectory = filepath.FromSlash("api/internal/db/migrations")
	migrationSource    = "file://" + migrationDirectory
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage DB migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up [num_steps]",
	Short: "Apply up migrations (defaults to running all migrations)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		m := getMigrator()

		if len(args) < 1 {
			guard(m.Up(), "run up migration")

			return
		}

		steps, err := strconv.ParseInt(args[0], 10, 32)
		guard(err, "parse number of migration steps")
		guard(m.Steps(int(steps)), "run up migration")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down [num_steps]",
	Short: "Apply down migrations (defaults to running all migrations)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		m := getMigrator()

		if len(args) < 1 {
			guard(m.Down(), "run up migration")

			return
		}

		steps, err := strconv.ParseInt(args[0], 10, 32)
		guard(err, "parse number of migration steps")

		if steps > 0 {
			steps = -steps
		}

		guard(m.Steps(int(steps)), "run up migration")
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create {migration_name}",
	Short: "Create new migration files",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		// Creating migration files isn't supported via the programmatic API, so create using the following shell command: `migrate create -seq -ext 'sql' -dir "api/internal/db/migrations" "$1"`
		migrator := exec.Command("migrate",
			"create",
			"-seq",
			"-ext", "sql",
			"-dir", migrationDirectory,
			args[0],
		)

		var migratorContent bytes.Buffer
		migrator.Stdout = os.Stdout
		migrator.Stderr = io.MultiWriter(os.Stderr, &migratorContent)
		guard(migrator.Run(), "execute `migrate ...` command")

		outLines := strings.Split(migratorContent.String(), "\n")
		for _, name := range outLines[:len(outLines)-1] {
			f, err := os.OpenFile(name, os.O_RDWR, 0o644)
			guard(err, "open migration file to add boilerplate")
			defer f.Close()
			_, err = f.WriteString("BEGIN;\n\n-- Migration logic goes here...\n\nCOMMIT;\n")
			guard(err, "write boilerplate to migration file")
		}
	},
}

var migrateFixCmd = &cobra.Command{
	Use:   "fix {version}",
	Short: "Reset the migration state of a DB to a known valid migration version, assuming all migrations were transacted",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		dbURL := getDatabaseURL(config.DBDriver, config.ServerBinName, config.DopplerEnvName, config.EnvVarPrefix, false)
		db, err := sql.Open(config.DBDriver.driverString(false), dbURL)
		guard(err, "open DB connection")

		version, err := strconv.ParseInt(args[0], 10, 32)
		guard(err, "parse desired version")

		if version < 1 {
			_, err = db.Exec("DROP TABLE IF EXISTS schema_migrations;")
			guard(err, "drop schema_migrations table")

			return
		}

		_, err = db.Exec("UPDATE schema_migrations SET version = $1, dirty = false;", version)
		guard(err, "reset schema_migrations table")
	},
}

func getMigrator() *migrate.Migrate {
	dbURL := getDatabaseURL(config.DBDriver, config.ServerBinName, config.DopplerEnvName, config.EnvVarPrefix, true)
	m, err := migrate.New(migrationSource, dbURL)
	guard(err, "construct DB migrator")

	return m
}
