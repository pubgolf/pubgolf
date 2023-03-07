package cmd

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
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
func (d DBDriver) driverString() string {
	switch d {
	case PostgreSQL:
		return "postgres"
	case SQLite3:
		return "sqlite3"
	}

	return ""
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
		c.ServerBinName = c.ProjectName + "-app-server"
	}

	if c.DopplerEnvName == "" {
		c.DopplerEnvName = "dev"
	}
}

// getDatabaseURL queries Doppler for DB connection details.
func getDatabaseURL(driver DBDriver, project, env, prefix string) string {
	vars := readDopplerVars(project, env, prefix, []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_SSL_MODE",
	})

	u := url.URL{
		User: url.UserPassword(
			getStr(vars, "DB_USER", ""),
			getStr(vars, "DB_PASSWORD", ""),
		),
		Host: getStr(vars, "DB_HOST", "localhost") + ":" + getStr(vars, "DB_PORT", "5432"),
		Path: getStr(vars, "DB_NAME", ""),
	}

	u.Scheme = driver.driverString()

	if driver == SQLite3 {
		u.Query().Set("x-no-tx-wrap", "true")
	}

	if driver == PostgreSQL {
		u.Query().Set("sslmoode", getStr(vars, "DB_SSL_MODE", "disable"))
	}

	return u.String()
}

// readDopplerVars queries the Doppler CLI for a requested set of computed env vars.
func readDopplerVars(project, env, prefix string, vars []string) map[string]string {
	doppler := exec.Command("doppler",
		"secrets",
		"--project", project,
		"--config", env,
		"--json",
	)

	var dopplerContent bytes.Buffer
	doppler.Stdout = &dopplerContent
	doppler.Stderr = os.Stderr

	guard(doppler.Run(), "execute `doppler ...` command")

	var data map[string]interface{}
	guard(json.NewDecoder(&dopplerContent).Decode(&data), "read JSON output from doppler")

	outData := make(map[string]string)
	for _, key := range vars {
		secret, ok := data[prefix+key]
		if !ok {
			continue
		}

		inner, ok := secret.(map[string]interface{})
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

		outData[key] = v
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