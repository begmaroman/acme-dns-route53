package acmstore

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/go-acme/lego/certificate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore"
	"github.com/begmaroman/acme-dns-route53/utils/strsl"
)

var (
	// ErrCertificateMissing is the error when certificate is empty
	ErrCertificateMissing = errors.New("certificate is empty")
)

// To make sure that ACM implements CertStore interface
var _ certstore.CertStore = &acmStore{}

// ACM is the implementation of CertStore interface.
// Used Amazon Certificate Manager to work with certificates
type acmStore struct {
	acm *acm.ACM
	log *logrus.Logger
}

// New is the constructor of acmStore
func New(provider client.ConfigProvider, log *logrus.Logger) certstore.CertStore {
	return &acmStore{
		acm: acm.New(provider),
		log: log,
	}
}

// Store implements CertStore interface
func (a *acmStore) Store(cert *certificate.Resource, domains []string) error {
	if cert == nil || cert.Certificate == nil {
		return ErrCertificateMissing
	}

	domainsListString := strings.Join(domains, ", ")

	a.log.Infof("[%s] acm: Retrieving server certificate", domainsListString)

	serverCert, err := retrieveServerCertificate(cert.Certificate)
	if err != nil {
		return errors.Wrap(err, "unable to retrieve server certificate")
	}

	a.log.Infof("[%s] acm: Finding existing server certificate in ACM", domainsListString)

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
		a.log.Infof("[%s] acm: Found existing server certificate in ACM with Arn = '%s'", domainsListString, aws.StringValue(certArn))
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

	a.log.Infof("[%s] acm: Imported certificate data in ACM with Arn = '%s'", domainsListString, aws.StringValue(resp.CertificateArn))

	return nil
}

// Load loads certificate by the given domains
func (a *acmStore) Load(domains []string) (*certstore.CertificateDetails, error) {
	cert, err := a.findExistingCertificate(domains)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find certificate")
	}

	return toCertificateDetails(cert), nil
}

// findExistingCertificate look ups a certificate in ACm by the given domains
func (a *acmStore) findExistingCertificate(domains []string) (*acm.CertificateDetail, error) {
	listResp, err := a.acm.ListCertificates(&acm.ListCertificatesInput{
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
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
