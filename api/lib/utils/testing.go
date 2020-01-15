package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	// Used to load the migrations files.
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/joho/godotenv"
)

func configureEnvVars() {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("error loading .env file", err)
	}
	os.Setenv("PUBGOLF_ENV", "test")
}

func getTestDB() *sql.DB {
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	if password != "" {
		password = ":" + password
	}
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("HOST_TEST_DB_PORT")
	if port == "" {
		port = "5432"
	}
	database := os.Getenv("TEST_DB_NAME")
	if database == "" {
		database = user
	}
	sslmode := os.Getenv("TEST_DB_SSL_MODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	connStr := fmt.Sprintf("postgres://%s%s@%s:%s/%s?sslmode=%s", user, password,
		host, port, database, sslmode)

	testDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("unable to connect to DB: %s", err)
	}

	return testDB
}

func resetTestDB(testDB *sql.DB) {
	driver, err := postgres.WithInstance(testDB, &postgres.Config{})
	migrator, err := migrate.NewWithDatabaseInstance("file://../../db/migrations", "postgres", driver)
	if err != nil {
		log.Fatal("error migrating DB", err)
	}

	migrator.Drop()
	migrator.Up()
}

// TestMain is an override for the per-file TestMain that sets up a test database and runs all seed statements before
// the unit tests.
func TestMain(m *testing.M, testDBHolder **sql.DB, seedStatements []string) {
	configureEnvVars()

	defer (*testDBHolder).Close()
	*testDBHolder = getTestDB()

	resetTestDB(*testDBHolder)

	for _, seedStatement := range seedStatements {
		_, err := (*testDBHolder).Exec(seedStatement)
		if err != nil {
			log.Fatal("error seeding DB", err)
		}
	}

	code := m.Run()

	os.Exit(code)
}
