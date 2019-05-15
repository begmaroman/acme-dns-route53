package lambda

import (
	"os"
	"strconv"
	"strings"
)

const (
	// DefaultRenewBefore is the default value of RENEW_BEFORE env var
	DefaultRenewBefore = 30
)

const (
	// DomainsEnvVar is the name of env var which contains domains list
	DomainsEnvVar = "DOMAINS"

	// LetsEncryptEnvVar is the name of env var which contains Let's Encrypt expiration email
	LetsEncryptEnvVar = "LETSENCRYPT_EMAIL"

	// StagingEnvVar is the name of env var which contains 1 value for using staging Let’s Encrypt environment or 0 for production environment.
	StagingEnvVar = "STAGING"

	// TopicEnvVar is the name of env var which contains a topic for notification
	TopicEnvVar = "NOTIFICATION_TOPIC"

	// RenewBeforeEnvVar is the name of env var which contains the number of days defining the period before expiration within which a certificate must be renewed
	RenewBeforeEnvVar = "RENEW_BEFORE"
)

// Config contains configuration data
type Config struct {
	Domains     []string
	Email       string
	Staging     bool
	Topic       string
	RenewBefore int
}

// InitConfig initializes configuration of the lambda function
func InitConfig(payload Payload) *Config {
	renewBefore, err := strconv.Atoi(os.Getenv(RenewBeforeEnvVar))
	if err != nil {
		renewBefore = DefaultRenewBefore
	}

	config := &Config{
		Domains:     strings.Split(os.Getenv(DomainsEnvVar), ","),
		Email:       os.Getenv(LetsEncryptEnvVar),
		Staging:     isStaging(os.Getenv(StagingEnvVar)),
		Topic:       os.Getenv(TopicEnvVar),
		RenewBefore: renewBefore,
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

	// Load notification topic
	if len(payload.Topic) > 0 {
		config.Topic = payload.Topic
	}

	// Load renew before days value
	if payload.RenewBefore > 0 {
		config.RenewBefore = payload.RenewBefore
	}

	return config
}

func isStaging(val string) bool {
	return val == "1"
}
