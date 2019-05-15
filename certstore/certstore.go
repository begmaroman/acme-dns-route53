package certstore

import (
	"time"

	"github.com/go-acme/lego/certificate"
)

// CertStore represents the interface to CRUD certificates
type CertStore interface {
	// Store represents logic to store the given certificate for the given domains
	Store(certificate *certificate.Resource, domains []string) error

	Load(domains []string) (*CertificateDetails, error)
}

// CertificateDetails contains certificate details
// TODO: Add more fields
type CertificateDetails struct {
	// NotAfter is the time after which the certificate is not valid.
	NotAfter time.Time
}
