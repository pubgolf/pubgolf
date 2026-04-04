// Package config provides a typed/structured config holder as well as logic to parse from env vars.
package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const envVarPrefix = "pubgolf"

// Init checks for env vars matching the `envVarPrefix` and returns a populated config struct.
func Init() (*App, error) {
	var c App

	err := envconfig.Process(envVarPrefix, &c)
	if err != nil {
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
	Port               int       `default:"5000"`
	EnvName            DeployEnv `required:"true" split_words:"true"`
	WebAppUpstreamHost string    `split_words:"true"`
	HostOrigin         string    `required:"true" split_words:"true"`

	// Behavioral config
	SMSAllowList PhoneNumSet `split_words:"true"`

	// Credentials
	HoneycombKey   string        `split_words:"true"`
	AppDatabaseURL string        `required:"true" split_words:"true"`
	Twilio         TwilioAuth    `required:"true" split_words:"true"`
	BlobStore      BlobStoreAuth `required:"true" split_words:"true"`

	// 1st party credentials and entropy
	AdminAuth WebAppAuth `required:"true" split_words:"true"`
}

// TwilioAuth contains Twilio credentials.
type TwilioAuth struct {
	AccountSID      string `required:"true" split_words:"true"`
	AuthToken       string `required:"true" split_words:"true"`
	VerificationSID string `required:"true" split_words:"true"`
}

// BlobStoreAuth contains credentials for the S3-compatible blob storage client.
type BlobStoreAuth struct {
	Endpoint  string `split_words:"true" required:"true"`
	AccessKey string `split_words:"true" required:"true"`
	SecretKey string `split_words:"true" required:"true"`
	Bucket    string `split_words:"true" required:"true"`
	UseSSL    bool   `split_words:"true" default:"false"`
}

// WebAppAuth contains auth params for the admin user.
type WebAppAuth struct {
	Password          string `required:"true" split_words:"true"`
	CookieToken       string `required:"true" split_words:"true"`
	AdminServiceToken string `required:"true" split_words:"true"`
}
