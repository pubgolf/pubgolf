package config

import (
	"errors"
	"fmt"
	"strings"
)

// DeployEnv is an enum allows toggling behavior based on the deployment environment.
type DeployEnv string

// Values for DeployEnv.
const (
	DeployEnvDev       DeployEnv = "development"
	DeployEnvE2ETest   DeployEnv = "e2e_test"
	DeployEnvPRPreview DeployEnv = "pr_preview"
	DeployEnvStaging   DeployEnv = "staging"
	DeployEnvProd      DeployEnv = "production"
)

// allDeployEnvs lists all valid values for a `DeployEnv` enum.
func allDeployEnvs() []DeployEnv {
	return []DeployEnv{
		DeployEnvDev,
		DeployEnvE2ETest,
		DeployEnvPRPreview,
		DeployEnvStaging,
		DeployEnvProd,
	}
}

// errInvalidDeployEnvValue indicates a parse error due to invalid enum value.
var errInvalidDeployEnvValue = errors.New("invalid enum value")

// Set attempts to parse a `DeployEnv` value from a string and returns an error on invalid values.
func (env *DeployEnv) Set(value string) error {
	v := strings.ToLower(value)
	for _, e := range allDeployEnvs() {
		if v == string(e) {
			*env = DeployEnv(v)

			return nil
		}
	}

	return fmt.Errorf("unrecognized value '%s': %w", v, errInvalidDeployEnvValue)
}
