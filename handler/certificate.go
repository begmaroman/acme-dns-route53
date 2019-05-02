package handler

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"

	"github.com/begmaroman/acme-dns-route53/certstore"
)

// CertificateHandler is the certificates handler
type CertificateHandler struct {
	isStaging bool
	configDir string

	store certstore.CertStore
	r53   *route53.Route53
}

// NewCertificateHandler is the constructor of CertificateHandler
func NewCertificateHandler(isStaging bool, store certstore.CertStore, r53 *route53.Route53, configDir string) *CertificateHandler {
	return &CertificateHandler{
		isStaging: isStaging,
		store:     store,
		r53:       r53,
		configDir: configDir,
	}
}

// toConfigParams creates a new configParams model
func (h *CertificateHandler) toConfigParams(user registration.User) *configParams {
	return &configParams{
		user:      user,
		isStaging: h.isStaging,
		keyType:   certcrypto.RSA2048, // TODO: Create a flag to define key type
	}
}

// toUserParams creates a new userParams model
func (h *CertificateHandler) toUserParams(email string) *userParams {
	return &userParams{
		email:   email,
		keyType: certcrypto.RSA2048, // TODO: Create a flag to define key type
	}
}
