package handler

import (
	"crypto"
	"encoding/pem"
	"os"

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"
	"github.com/pkg/errors"
)

// CertUser is the simple implementation of acme.User interface
type CertUser struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	key          crypto.PrivateKey
}

// NewCertUser is the constructor of CertUser
func NewCertUser(email string) *CertUser {
	return &CertUser{
		Email: email,
	}
}

// GetEmail returns email of the user
func (u *CertUser) GetEmail() string {
	return u.Email
}

// GetRegistration returns registration.Resource model of the user
func (u CertUser) GetRegistration() *registration.Resource {
	return u.Registration
}

// GetPrivateKey returns the private key of the user
func (u *CertUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// StorePrivateKey stores the private key to the file by the given path
// configDir - is the root of configs. Must be present without "/" in the end.
// TODO: Create an interface for storing user's private key
func (u *CertUser) StorePrivateKey(configDir string) error {
	if len(configDir) > 0 {
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			if err := os.MkdirAll(configDir, 0666); err != nil {
				return errors.Wrap(err, "unable to create config directory")
			}
		}

		configDir += "/"
	}

	filePath := configDir + u.Email + ".pem"

	certOut, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "unable to create file with path '%s'", filePath)
	}
	defer certOut.Close()

	pemKey := certcrypto.PEMBlock(u.key)
	if err = pem.Encode(certOut, pemKey); err != nil {
		return errors.Wrap(err, "unable to encode private key")
	}

	return nil
}
