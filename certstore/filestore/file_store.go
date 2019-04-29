package filestore

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/go-acme/lego/certificate"
	"github.com/pkg/errors"

	"github.com/begmaroman/acme-dns-route53/certstore"
)

var (
	ErrPrivateKeyMissing = errors.New("unable to save pem without private key")
)

// Ensures that FileStore implements CertStore interface
var _ certstore.CertStore = &FileStore{}

// FileStore is the implementation of CertStore interface.
// Used file system to work with certificates
type FileStore struct {
}

// New is the constructor of FileStore
func New() *FileStore {
	return &FileStore{}
}

// Store implements CertStore interface
func (f *FileStore) Store(cert *certificate.Resource) error {
	domain := cert.Domain
	pem := true // TODO: Fix this

	// We store the certificate, private key and metadata in different files
	// as web servers would not be able to work with a combined file.
	if err := ioutil.WriteFile(domain+".crt", cert.Certificate, 0600); err != nil {
		return errors.Wrapf(err, "unable to save Certificate for domain %s", domain)
	}

	if cert.IssuerCertificate != nil {
		if err := ioutil.WriteFile(domain+".issuer.crt", cert.IssuerCertificate, 0600); err != nil {
			return errors.Wrapf(err, "unable to save IssuerCertificate for domain %s", domain)
		}
	}

	if cert.PrivateKey != nil {
		// if we were given a CSR, we don't know the private key
		if err := ioutil.WriteFile(domain+".key", cert.PrivateKey, 0600); err != nil {
			return errors.Wrapf(err, "unable to save PrivateKey for domain %s", domain)
		}

		if pem {
			if err := ioutil.WriteFile(domain+".pem", bytes.Join([][]byte{cert.Certificate, cert.PrivateKey}, nil), 0600); err != nil {
				return errors.Wrapf(err, "unable to save Certificate and PrivateKey in .pem for domain %s", domain)
			}
		}
	} else if pem {
		// We don't have the private key; can't write the .pem file
		return ErrPrivateKeyMissing
	}

	jsonBytes, err := json.MarshalIndent(cert, "", "\t")
	if err != nil {
		return errors.Wrapf(err, "unable to marshal CertResource for domain %s", domain)
	}

	if err = ioutil.WriteFile(domain+".json", jsonBytes, 0600); err != nil {
		return errors.Wrapf(err, "unable to save CertResource for domain %s", domain)
	}

	return nil
}
