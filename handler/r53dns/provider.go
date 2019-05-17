package r53dns

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// dnsProvider is the custom implementation of the challenge.Provider interface
type dnsProvider struct {
	r53Worker *r53ResourceWorker
	log       *logrus.Logger
}

// New is the constructor of DNSProvider
func New(provider client.ConfigProvider, log *logrus.Logger) challenge.Provider {
	return &dnsProvider{
		r53Worker: newR53ResourceWorker(route53.New(provider), log),
		log:       log,
	}
}

// Present prints instructions for manually creating the TXT record
func (p *dnsProvider) Present(domain, token, keyAuth string) error {
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
func (p *dnsProvider) CleanUp(domain, token, keyAuth string) error {
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
