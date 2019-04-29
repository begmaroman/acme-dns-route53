package acmstore

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/go-acme/lego/certificate"
	"github.com/pkg/errors"

	"github.com/begmaroman/acme-dns-route53/certstore"
)

var (
	// ErrCertificateMissing is the error when certificate is empty
	ErrCertificateMissing = errors.New("certificate is empty")
)

// Ensures that ACM implements CertStore interface
var _ certstore.CertStore = &ACM{}

// ACM is the implementation of CertStore interface.
// Used Amazon Certificate Manager to work with certificates
type ACM struct {
	acm *acm.ACM
}

// New is the constructor of ACM
func New(provider client.ConfigProvider) *ACM {
	return &ACM{
		acm: acm.New(provider),
	}
}

// Store implements CertStore interface
func (a *ACM) Store(cert *certificate.Resource) error {
	if cert == nil || cert.Certificate == nil {
		return ErrCertificateMissing
	}

	serverCert, err := retrieveServerCertificate(cert.Certificate)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve server certificate")
	}

	// Init request parameters
	input := &acm.ImportCertificateInput{
		Certificate: serverCert,
		PrivateKey:  cert.PrivateKey,
	}

	resp, err := a.acm.ImportCertificate(input)
	if err != nil {
		return errors.Wrap(err, "unable to store certificate into ACM")
	}

	fmt.Println("resp.CertificateArn", resp.CertificateArn)

	return nil
}

// retrieveServerCertificate retrieves the server certificate from the given PEM encoded list
func retrieveServerCertificate(list []byte) ([]byte, error) {
	var blocks []*pem.Block
	for {
		var certDERBlock *pem.Block
		certDERBlock, list = pem.Decode(list)
		if certDERBlock == nil {
			break
		}

		if certDERBlock.Type == "CERTIFICATE" {
			blocks = append(blocks, certDERBlock)
		}
	}

	crt := bytes.NewBuffer(nil)
	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse certificate")
		}

		if !cert.IsCA {
			pem.Encode(crt, block)
			break
		}
	}

	return crt.Bytes(), nil
}
