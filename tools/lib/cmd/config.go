package cmd

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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
		return "pgx"
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
		c.ServerBinName = c.ProjectName + "-api-server"
	}

	if c.DopplerEnvName == "" {
		c.DopplerEnvName = "dev"
	}
}

// getDatabaseURL queries Doppler for DB connection details.
func getDatabaseURL(driver DBDriver, project, env, prefix string) string {
	vars := readDopplerVars(project, env, prefix, []string{
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

	u := url.URL{
		Scheme: driver.driverString(),
	}

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
