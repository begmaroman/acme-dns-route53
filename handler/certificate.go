package handler

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/registration"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore"
)

// CertificateHandlerOptions is the options of certificate handler
type CertificateHandlerOptions struct {
	Staging   bool
	ConfigDir string
	Store     certstore.CertStore
	R53       *route53.Route53
	SNS       *sns.SNS
	Log       *logrus.Logger
}

// CertificateHandler is the certificates handler
type CertificateHandler struct {
	isStaging bool
	configDir string

	store certstore.CertStore
	sns   *sns.SNS
	r53   *route53.Route53
	log   *logrus.Logger
}

// NewCertificateHandler is the constructor of CertificateHandler
func NewCertificateHandler(opts *CertificateHandlerOptions) *CertificateHandler {
	return &CertificateHandler{
		isStaging: opts.Staging,
		store:     opts.Store,
		sns:       opts.SNS,
		r53:       opts.R53,
		configDir: opts.ConfigDir,
		log:       opts.Log,
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
