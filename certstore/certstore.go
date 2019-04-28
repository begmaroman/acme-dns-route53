package certstore

import "github.com/go-acme/lego/certificate"

// CertStore represents the interface to CRUD certificates
type CertStore interface {
	Store(certificate *certificate.Resource) error
}
