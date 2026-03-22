package cmd

import (
	"context"
	"net/url"
	"path/filepath"
	"strings"
)

// DBDriver is an enum specifying which database driver to use for migrations and DAO codegen.
type DBDriver int

// DBDriver values.
const (
	None DBDriver = iota
	PostgreSQL
	SQLite3
)

// driverString returns the database/sql driver ID for the given database type.
func (d DBDriver) driverString(isMigrator bool) string {
	switch d {
	case PostgreSQL:
		// Hack to deal with golang-migrate needing a different driver string for pgx vs pq.
		if isMigrator {
			return "pgx5"
		}

		return "pgx"
	case SQLite3:
		return "sqlite3"
	case None:
		return ""
	default:
		return ""
	}
}

// CLIConfig sets naming and capabilities for the generated CLI tool.
type CLIConfig struct {
	ProjectName    string
	CLIName        string
	ServerBinName  string
	DopplerEnvName string
	EnvVarPrefix   string
	DBDriver       DBDriver
}

func (c *CLIConfig) setDefaults() {
	if c.CLIName == "" {
		c.CLIName = c.ProjectName + "-devctrl"
	}

	if c.ServerBinName == "" {
		c.ServerBinName = c.ProjectName + "-api-server"
	}

	if c.DopplerEnvName == "" {
		c.DopplerEnvName = "dev"
	}
}

// getDatabaseURL uses the EnvProvider to fetch DB connection details and build a URL.
func getDatabaseURL(ctx context.Context, ep EnvProvider, driver DBDriver, project, env, prefix string, isMigrator bool) string {
	vars := readEnvVars(ctx, ep, project, env, prefix, []string{
		"SQLITE_PATH",
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_SSL_MODE",
	})

	if driver == SQLite3 {
		return "sqlite3://" + getStr(vars, "SQLITE_PATH", filepath.FromSlash("./data/db/data.db")) + "?x-no-tx-wrap=true"
	}

	scheme := "postgres"
	if isMigrator {
		scheme = "pgx5"
	}

	u := url.URL{Scheme: scheme}
	u.User = url.UserPassword(
		getStr(vars, "DB_USER", config.ProjectName+"_dev"),
		getStr(vars, "DB_PASSWORD", config.ProjectName+"_dev"),
	)
	u.Host = getStr(vars, "DB_HOST", "localhost") + ":" + getStr(vars, "DB_PORT", "5432")
	u.Path = getStr(vars, "DB_NAME", config.ProjectName+"_dev")

	q := u.Query()
	q.Add("sslmode", getStr(vars, "DB_SSL_MODE", "disable"))
	u.RawQuery = q.Encode()

	return u.String()
}

// readEnvVars fetches env vars from the provider and extracts the requested
// keys (with the configured prefix stripped).
func readEnvVars(ctx context.Context, ep EnvProvider, project, env, prefix string, vars []string) map[string]string {
	envSlice, err := ep.Env(ctx, project, env)
	guard(err, "fetch environment variables")

	// Build a lookup from the returned KEY=VALUE pairs.
	envMap := make(map[string]string, len(envSlice))
	for _, entry := range envSlice {
		k, v, _ := strings.Cut(entry, "=")
		envMap[k] = v
	}

	// Extract requested vars, stripping the prefix.
	outData := make(map[string]string)

	for _, key := range vars {
		if v, ok := envMap[prefix+key]; ok {
			outData[key] = v
		}
	}

	return outData
}

// getStr is a map getter with a default value.
func getStr(m map[string]string, k, def string) string {
	if v, ok := m[k]; ok {
		return v
	}

	return def
}
