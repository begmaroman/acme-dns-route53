package handler

import (
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/log"
	"github.com/go-acme/lego/registration"
	"github.com/pkg/errors"
)

// configParams is the parameters which are needed for config creation
type configParams struct {
	isStaging bool
	keyType   certcrypto.KeyType
	user      registration.User
}

// getConfig creates a config for the lego client
func getConfig(params *configParams) (*lego.Config, error) {
	// Create a new config
	config := lego.NewConfig(params.user)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	if params.isStaging {
		log.Infof("acme: Using staging environment")
		config.CADirURL = lego.LEDirectoryStaging
	}
	config.Certificate.KeyType = params.keyType

	return config, nil
}

// configParams is the parameters which are needed for user loading
type userParams struct {
	email   string
	keyType certcrypto.KeyType
}

// getUser create a new getUser object
// TODO: Load user based on a data in config directory
func getUser(params *userParams) (*CertUser, error) {
	var err error

	// Create a user
	certUser := NewCertUser(params.email)

	// New accounts need a private key to start
	if certUser.key, err = certcrypto.GeneratePrivateKey(params.keyType); err != nil {
		return nil, errors.Wrap(err, "unable to generate private key")
	}

	return certUser, nil
}
