package handler

import (
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/begmaroman/acme-dns-route53/certstore"
	"github.com/begmaroman/acme-dns-route53/printer"
)

// CertificateHandler is the certificates handler
type CertificateHandler struct {
	printer printer.Printer
	store   certstore.CertStore
	r53     *route53.Route53
}

// NewCertificateHandler is the constructor of CertificateHandler
func NewCertificateHandler(printer printer.Printer, store certstore.CertStore, r53 *route53.Route53) *CertificateHandler {
	return &CertificateHandler{
		printer: printer,
		store:   store,
		r53:     r53,
	}
}
