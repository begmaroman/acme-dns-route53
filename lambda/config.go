package lambda

import (
	"os"
	"strings"
)

const (
	// DomainsEnvVar is the name of env var which contains domains list
	DomainsEnvVar = "DOMAINS"

	// LetsEncryptEnvVar is the name of env var which contains Let's Encrypt expiration email
	LetsEncryptEnvVar = "LETSENCRYPT_EMAIL"

	// StagingEnvVar is the name of env var which contains 1 value for using staging Let’s Encrypt environment or 0 for production environment.
	StagingEnvVar = "STAGING"
)

// Config contains configuration data
type Config struct {
	Domains []string
	Email   string
	Staging bool
}

// InitConfig initializes configuration of the lambda function
func InitConfig(payload Payload) *Config {
	config := &Config{
		Domains: strings.Split(os.Getenv(DomainsEnvVar), ","),
		Email:   os.Getenv(LetsEncryptEnvVar),
		Staging: isStaging(os.Getenv(StagingEnvVar)),
	}

	// Load domains
	if len(payload.Domains) > 0 {
		config.Domains = payload.Domains
	}

	// Load email
	if len(payload.Email) > 0 {
		config.Email = payload.Email
	}

	// Load environment
	if len(payload.Staging) > 0 {
		config.Staging = isStaging(payload.Staging)
	}

	return config
}

func isStaging(val string) bool {
	return val == "1"
}
