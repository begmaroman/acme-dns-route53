package handler

import (
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"

	"github.com/begmaroman/acme-dns-route53/handler/r53dns"
)

var (
	// registerOptions is the predefined registration.RegisterOptions struct with the default params
	registerOptions = registration.RegisterOptions{TermsOfServiceAgreed: true}
)

// Create creates new SSL certificates for the given domains with the given email
func (h *CertificateHandler) Create(domains []string, email string) error {
	var err error

	// Create a user
	certUser := NewCertUser(email)

	// New accounts need a private key to start
	// TODO: Maybe we need to store user's private key
	if certUser.key, err = certcrypto.GeneratePrivateKey(certcrypto.RSA2048); err != nil {
		return err
	}

	// Create a new config
	config := lego.NewConfig(certUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = lego.LEDirectoryStaging       // TODO: Create a fleg to define production or staging
	config.Certificate.KeyType = certcrypto.RSA2048 // TODO: Create a flag to define key type

	// Create a client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	// Use DNS-01 challenge to verify that the given domain belongs to the current server
	if err = client.Challenge.SetDNS01Provider(r53dns.NewProvider(h.r53)); err != nil {
		return err
	}

	// New users will need to register
	if certUser.Registration, err = client.Registration.Register(registerOptions); err != nil {
		return err
	}

	// Create a new request to obtain certificate
	request := certificate.ObtainRequest{
		Domains:    domains,
		PrivateKey: nil,
		Bundle:     true,
		MustStaple: false,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	// Store the obtained certificate
	if err := h.store.Store(certificates); err != nil {
		return err
	}

	// TODO: Do smth with certUser

	return nil
}
