package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// getPostgresURL queries Doppler for DB connection details.
func getPostgresURL(project, env, prefix string) string {
	vars := readDopplerVars(project, env, prefix, []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	})

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getStr(vars, "DB_USER", ""),
		getStr(vars, "DB_PASSWORD", ""),
		getStr(vars, "DB_HOST", "localhost"),
		getStr(vars, "DB_PORT", "5432"),
		getStr(vars, "DB_NAME", ""),
		getStr(vars, "DB_SSL_MODE", "disable"),
	)
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
