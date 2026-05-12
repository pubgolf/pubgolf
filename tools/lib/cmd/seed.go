package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pubgolf/pubgolf/api/db/seeds"

	// pgx driver for database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	errSeedEventExists  = errors.New("event already exists; use `seed reset` to recreate")
	errSeedNotFound     = errors.New("unknown seed name")
	errSeedMissingDBURL = errors.New("APP_DATABASE_URL not found in Doppler config")
)

func init() {
	seedCmd.AddCommand(seedCreateCmd)
	seedCmd.AddCommand(seedResetCmd)
	seedCmd.AddCommand(seedCleanCmd)
	seedCmd.AddCommand(seedStatusCmd)
	seedCmd.PersistentFlags().String("env", "", "Target environment (e.g. stg, staging, prd, prod)")
	rootCmd.AddCommand(seedCmd)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Manage seed data for development and testing",
	Run: func(_ *cobra.Command, _ []string) {
		log.Println("Available seeds:")

		for _, name := range seeds.Names() {
			s := seeds.Registry[name]
			log.Printf("  %s (event: %s)", name, s.EventKey)
		}
	},
}

var seedCreateCmd = &cobra.Command{
	Use:   "create {seed_name}",
	Short: "Create seed data (fails if event already exists)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		classifyAndExit(seedCreate(cmd.Context(), envProvider, args[0], env))
	},
}

var seedResetCmd = &cobra.Command{
	Use:   "reset {seed_name}",
	Short: "Delete and recreate seed data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		classifyAndExit(seedReset(cmd.Context(), envProvider, args[0], env, dryRun))
	},
}

var seedCleanCmd = &cobra.Command{
	Use:   "clean {seed_name}",
	Short: "Delete seed data",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		classifyAndExit(seedClean(cmd.Context(), envProvider, args[0], env, dryRun))
	},
}

var seedStatusCmd = &cobra.Command{
	Use:   "status {seed_name}",
	Short: "Show seed data status and detect drift from expected baseline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, _ := cmd.Flags().GetString("env")
		classifyAndExit(seedStatus(cmd.Context(), envProvider, args[0], env))
	},
}

// seedCreate inserts seed data, failing if the event already exists.
func seedCreate(ctx context.Context, ep EnvProvider, seedName, envName string) error {
	s, err := lookupSeed(seedName)
	if err != nil {
		return err
	}

	db, err := openSeedDB(ctx, ep, envName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, exists, err := seeds.PreviewEventData(ctx, db, s.EventKey)
	if err != nil {
		return fmtErr(err, "check existing event data")
	}

	if exists {
		return fmt.Errorf("event %q: %w", s.EventKey, errSeedEventExists)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmtErr(err, "begin transaction")
	}
	defer tx.Rollback() //nolint:errcheck

	runErr := s.Run(ctx, tx)
	if runErr != nil {
		return fmtErr(runErr, "run seed")
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return fmtErr(commitErr, "commit seed transaction")
	}

	log.Printf("Created seed %q (event: %s)\n", s.Name, s.EventKey)

	return printSeedStatus(ctx, db, s)
}

// seedReset deletes existing data and recreates the seed.
func seedReset(ctx context.Context, ep EnvProvider, seedName, envName string, dryRun bool) error {
	s, err := lookupSeed(seedName)
	if err != nil {
		return err
	}

	db, err := openSeedDB(ctx, ep, envName)
	if err != nil {
		return err
	}
	defer db.Close()

	summary, exists, err := seeds.PreviewEventData(ctx, db, s.EventKey)
	if err != nil {
		return fmtErr(err, "check existing event data")
	}

	if dryRun {
		if exists {
			log.Println("Dry run: would delete existing data:")
			printEventSummary(summary)
		} else {
			log.Println("Dry run: no existing data to delete")
		}

		log.Println("Dry run: would then create seed data")

		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmtErr(err, "begin transaction")
	}
	defer tx.Rollback() //nolint:errcheck

	if exists {
		delErr := seeds.DeleteEventData(ctx, tx, s.EventKey)
		if delErr != nil {
			return fmtErr(delErr, "delete existing event data")
		}

		log.Printf("Deleted existing data for event %q\n", s.EventKey)
	}

	runErr := s.Run(ctx, tx)
	if runErr != nil {
		return fmtErr(runErr, "run seed")
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return fmtErr(commitErr, "commit seed transaction")
	}

	log.Printf("Reset seed %q (event: %s)\n", s.Name, s.EventKey)

	return printSeedStatus(ctx, db, s)
}

// seedClean deletes all data for a seed's event.
func seedClean(ctx context.Context, ep EnvProvider, seedName, envName string, dryRun bool) error {
	s, err := lookupSeed(seedName)
	if err != nil {
		return err
	}

	db, err := openSeedDB(ctx, ep, envName)
	if err != nil {
		return err
	}
	defer db.Close()

	summary, exists, err := seeds.PreviewEventData(ctx, db, s.EventKey)
	if err != nil {
		return fmtErr(err, "check existing event data")
	}

	if !exists {
		log.Printf("No data found for event %q — nothing to clean\n", s.EventKey)

		return nil
	}

	if dryRun {
		log.Println("Dry run: would delete:")
		printEventSummary(summary)

		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmtErr(err, "begin transaction")
	}
	defer tx.Rollback() //nolint:errcheck

	delErr := seeds.DeleteEventData(ctx, tx, s.EventKey)
	if delErr != nil {
		return fmtErr(delErr, "delete event data")
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return fmtErr(commitErr, "commit clean transaction")
	}

	log.Printf("Cleaned seed %q (event: %s)\n", s.Name, s.EventKey)

	return nil
}

// seedStatus shows the current state of seed data vs expected baseline.
func seedStatus(ctx context.Context, ep EnvProvider, seedName, envName string) error {
	s, err := lookupSeed(seedName)
	if err != nil {
		return err
	}

	db, err := openSeedDB(ctx, ep, envName)
	if err != nil {
		return err
	}
	defer db.Close()

	return printSeedStatus(ctx, db, s)
}

// lookupSeed finds a seed by name in the registry.
func lookupSeed(name string) (seeds.Seed, error) {
	s, ok := seeds.Registry[name]
	if !ok {
		return seeds.Seed{}, fmt.Errorf("%w: %q; available: %s", errSeedNotFound, name, strings.Join(seeds.Names(), ", "))
	}

	return s, nil
}

// openSeedDB opens a database connection for the given environment.
func openSeedDB(ctx context.Context, ep EnvProvider, envName string) (*sql.DB, error) {
	var dbURL string

	if envName == "" {
		dbURL = getDatabaseURL(ctx, ep, config.DBDriver, config.ServerBinName, config.DopplerEnvName, config.EnvVarPrefix, false)
	} else {
		dopplerCfg := normalizeEnvName(envName)
		vars := readEnvVars(ctx, ep, config.ServerBinName, dopplerCfg, config.EnvVarPrefix, []string{"APP_DATABASE_URL"})

		var ok bool

		dbURL, ok = vars["APP_DATABASE_URL"]
		if !ok {
			return nil, fmt.Errorf("%w: %q", errSeedMissingDBURL, dopplerCfg)
		}
	}

	db, err := sql.Open(config.DBDriver.driverString(false), dbURL)
	if err != nil {
		return nil, fmtErr(err, "open database connection")
	}

	// Verify connectivity.
	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		db.Close()

		return nil, fmtErr(pingErr, "ping database")
	}

	return db, nil
}

// normalizeEnvName maps common environment name aliases to Doppler config names.
func normalizeEnvName(env string) string {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "stg", "staging", "stage":
		return "stg"
	case "prd", "prod", "production":
		return "prd"
	case "dev", "development":
		return "dev"
	default:
		return env // Pass through as-is; Doppler will error if invalid.
	}
}

// printSeedStatus queries the database and prints a comparison of actual vs expected counts.
func printSeedStatus(ctx context.Context, db *sql.DB, s seeds.Seed) error {
	summary, exists, err := seeds.PreviewEventData(ctx, db, s.EventKey)
	if err != nil {
		return fmtErr(err, "query event data")
	}

	log.Printf("Seed %q (event: %s)\n", s.Name, s.EventKey)

	if !exists {
		log.Println("  Event does not exist")

		return nil
	}

	type row struct {
		label    string
		expected int
		actual   int
	}

	rows := []row{
		{"Stages", s.Expected.Stages, summary.Stages},
		{"Rules", s.Expected.Rules, summary.Rules},
		{"Players", s.Expected.Players, summary.Players},
		{"Scores", s.Expected.Scores, summary.Scores},
		{"Adj. templates", s.Expected.AdjustmentTemplates, summary.AdjustmentTemplates},
	}

	log.Println()
	log.Printf("  %-20s %8s %8s", "Metric", "Expected", "Actual")
	log.Printf("  %-20s %8s %8s", "------", "--------", "------")

	for _, r := range rows {
		marker := "ok"

		if r.actual != r.expected {
			diff := r.actual - r.expected
			if diff > 0 {
				marker = fmt.Sprintf("+%d", diff)
			} else {
				marker = strconv.Itoa(diff)
			}
		}

		log.Printf("  %-20s %8d %8d  %s", r.label, r.expected, r.actual, marker)
	}

	if summary.Scores > 0 {
		log.Printf("\n  Scores: %d verified, %d unverified", summary.ScoresVerified, summary.ScoresUnverified)
	}

	if len(summary.PlayerNames) > 0 {
		log.Printf("  Sample players: %s", strings.Join(summary.PlayerNames, ", "))
	}

	log.Println()

	return nil
}

// printEventSummary logs the counts in an event summary.
func printEventSummary(s seeds.EventSummary) {
	log.Printf("  Stages:              %d", s.Stages)
	log.Printf("  Rules:               %d", s.Rules)
	log.Printf("  Players:             %d", s.Players)
	log.Printf("  Scores:              %d (%d verified, %d unverified)", s.Scores, s.ScoresVerified, s.ScoresUnverified)
	log.Printf("  Adj. templates:      %d", s.AdjustmentTemplates)

	if len(s.PlayerNames) > 0 {
		log.Printf("  Sample players:      %s", strings.Join(s.PlayerNames, ", "))
	}
}
