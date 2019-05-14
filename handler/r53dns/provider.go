package r53dns

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Ensures that DNSProvider implements challenge.Provider interface
var _ challenge.Provider = &DNSProvider{}

// DNSProvider is the custom implementation of the challenge.Provider interface
type DNSProvider struct {
	r53Worker *r53ResourceWorker
	log       *logrus.Logger
}

// NewDNSProviderManual is the constructor of DNSProvider
func NewProvider(r53 *route53.Route53, log *logrus.Logger) *DNSProvider {
	return &DNSProvider{
		r53Worker: newR53ResourceWorker(r53, log),
		log:       log,
	}
}

// Present prints instructions for manually creating the TXT record
func (p *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return errors.Wrapf(err, "unable to find zone by FQDN = '%s'", fqdn)
	}

	p.log.Infof("[%s] acme: Creating TXT record in %s zone", domain, authZone)

	// Create a subdomain
	recordID, err := p.r53Worker.changeDNSRecord(route53.ChangeActionUpsert, fqdn, buildQuotedValue(value))
	if err != nil {
		return errors.Wrapf(err, "unable to change a record with FQDN = '%s'", fqdn)
	}

	p.log.Infof("[%s] acme: Created TXT record in %s zone with ID %s", domain, authZone, recordID)

	return err
}

// CleanUp prints instructions for manually removing the TXT record
func (p *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)

	// Retrieve zone by FQDN
	authZone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return errors.Wrapf(err, "unable to find zone by FQDN = '%s'", fqdn)
	}

	p.log.Infof("[%s] acme: Removing TXT record from %s zone", domain, authZone)

	// Delete the subdomain
	recordID, err := p.r53Worker.changeDNSRecord(route53.ChangeActionDelete, fqdn, buildQuotedValue(value))
	if err != nil {
		return errors.Wrapf(err, "unable to delete a record with FQDN = '%s'", fqdn)
	}

	p.log.Infof("[%s] acme: Removed TXT record in %s zone with ID %s", domain, authZone, recordID)

	return nil
}
