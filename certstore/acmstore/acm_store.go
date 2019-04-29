package acmstore

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/log"
	"github.com/pkg/errors"

	"github.com/begmaroman/acme-dns-route53/certstore"
	"github.com/begmaroman/acme-dns-route53/utils/strsl"
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
		acm: acm.New(provider, aws.NewConfig().WithRegion("eu-central-1")),
	}
}

// Store implements CertStore interface
func (a *ACM) Store(cert *certificate.Resource, domains []string) error {
	if cert == nil || cert.Certificate == nil {
		return ErrCertificateMissing
	}

	domainsListString := strings.Join(domains, ", ")

	log.Infof("[%s] acm: Retrieving server certificate", domainsListString)

	serverCert, err := retrieveServerCertificate(cert.Certificate)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve server certificate")
	}

	log.Infof("[%s] acm: Finding existing server certificate in ACM", domainsListString)

	existingCert, err := a.findExistingCertificate(domains)
	if err != nil {
		return errors.Wrap(err, "unable to find existing certificate")
	}

	// Retrieve exising certificate ID
	var certArn *string
	if existingCert != nil {
		certArn = existingCert.CertificateArn
	}

	if certArn != nil {
		log.Infof("[%s] acm: Found existing server certificate in ACM with Arn = '%s'", domainsListString, certArn)
	}

	// Init request parameters
	input := &acm.ImportCertificateInput{
		Certificate:      serverCert,
		CertificateArn:   certArn,
		CertificateChain: cert.IssuerCertificate,
		PrivateKey:       cert.PrivateKey,
	}

	resp, err := a.acm.ImportCertificate(input)
	if err != nil {
		return errors.Wrap(err, "unable to store certificate into ACM")
	}

	log.Infof("[%s] acm: Imported certificate data in ACM with Arn = '%s'", domainsListString, *resp.CertificateArn)

	return nil
}

// findExistingCertificate look ups a certificate in ACm by the given domains
func (a *ACM) findExistingCertificate(domains []string) (*acm.CertificateDetail, error) {
	listResp, err := a.acm.ListCertificates(&acm.ListCertificatesInput{
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		fmt.Println("err", err)
		return nil, errors.Wrap(err, "unable to list certificates")
	}

	for _, crt := range listResp.CertificateSummaryList {
		certResp, err := a.acm.DescribeCertificate(&acm.DescribeCertificateInput{
			CertificateArn: crt.CertificateArn,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unable to describe certificate")
		}

		altNames := aws.StringValueSlice(certResp.Certificate.SubjectAlternativeNames)
		if strsl.ContainsSub(domains, altNames) {
			return certResp.Certificate, nil
		}
	}

	return nil, nil
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
