// Package testguard manages flags for use in tests.
package testguard

import (
	"flag"
	"os"
)

var (
	enableE2ETests    = flag.Bool("e2e", false, "run e2e tests")
	enablePostgresLog = flag.Bool("postgres-log", false, "emit logs for embedded postgres instance (only works when --shared-postgres=false)")
	useSharedPostgres = flag.Bool("shared-postgres", false, "run database tests against a running postgres instance, provided via PUBGOLF_SHARED_DB_URL env var")
)

// UnitTest marks this test suite as unit-level and skips running if `-e2e=true` is passed to the test command.
func UnitTest() {
	flag.Parse()

	if *enableE2ETests {
		os.Exit(0)
	}
}

// E2ETest marks this test suite as e2e-level and skips running if `-e2e=true` is not passed to the test command.
func E2ETest() {
	flag.Parse()

	if !*enableE2ETests {
		os.Exit(0)
	}
}

// DBFlags stores test preferences related to database config.
type DBFlags struct {
	UseSharedPostgres bool
	EnablePostgresLog bool
}

// GetDBFlags gets the set flag values for DB-related test flags.
func GetDBFlags() DBFlags {
	flag.Parse()

	return DBFlags{
		UseSharedPostgres: *useSharedPostgres,
		EnablePostgresLog: *enablePostgresLog,
	}
}
