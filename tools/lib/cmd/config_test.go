package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIConfig_SetDefaults(t *testing.T) {
	t.Parallel()

	c := CLIConfig{ProjectName: "pubgolf"}
	c.setDefaults()

	assert.Equal(t, "pubgolf-devctrl", c.CLIName)
	assert.Equal(t, "pubgolf-api-server", c.ServerBinName)
	assert.Equal(t, "dev", c.DopplerEnvName)
}

func TestCLIConfig_SetDefaults_PreservesExplicitValues(t *testing.T) {
	t.Parallel()

	c := CLIConfig{
		ProjectName:    "pubgolf",
		CLIName:        "custom-cli",
		ServerBinName:  "custom-server",
		DopplerEnvName: "staging",
	}
	c.setDefaults()

	assert.Equal(t, "custom-cli", c.CLIName)
	assert.Equal(t, "custom-server", c.ServerBinName)
	assert.Equal(t, "staging", c.DopplerEnvName)
}

func TestGetStr_ReturnsValue(t *testing.T) {
	t.Parallel()

	m := map[string]string{"key": "value"}
	assert.Equal(t, "value", getStr(m, "key", "default"))
}

func TestGetStr_ReturnsDefault(t *testing.T) {
	t.Parallel()

	m := map[string]string{}
	assert.Equal(t, "default", getStr(m, "missing", "default"))
}

func TestDBDriverString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		driver     DBDriver
		isMigrator bool
		want       string
	}{
		{name: "postgres non-migrator", driver: PostgreSQL, isMigrator: false, want: "pgx"},
		{name: "postgres migrator", driver: PostgreSQL, isMigrator: true, want: "pgx5"},
		{name: "sqlite3", driver: SQLite3, isMigrator: false, want: "sqlite3"},
		{name: "none", driver: None, isMigrator: false, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.driver.driverString(tt.isMigrator))
		})
	}
}
