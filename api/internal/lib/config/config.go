// Package config provides a typed/structured config holder.
package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const envVarPrefix = "pubgolf"

// Init checks for env vars matchint the `envVarPrefix` and returns a populated config struct.
func Init() (*App, error) {
	var c App
	if err := envconfig.Process(envVarPrefix, &c); err != nil {
		return nil, fmt.Errorf("read env vars into config: %w", err)
	}

	// Unfortunate Render-related hackery, since we can't override the `PUBGOLF_ENV_NAME` var directly in the preview env.
	if os.Getenv("IS_PULL_REQUEST") == "true" {
		c.EnvName = DeployEnvPRPreview
	}

	return &c, nil
}

// App defines the env config for the app. The tag directive syntax is defined at https://github.com/kelseyhightower/envconfig.
type App struct {
	// Env config
	Port    int       `default:"5000"`
	EnvName DeployEnv `required:"true" split_words:"true"`

	// Credentials
	HoneycombWriteKey string `split_words:"true"`
	AppDatabaseURL    string `required:"true" split_words:"true"`
}