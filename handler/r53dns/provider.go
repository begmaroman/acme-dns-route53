package r53dns

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/challenge"
	"github.com/go-acme/lego/challenge/dns01"
	"github.com/go-acme/lego/log"
	"github.com/pkg/errors"

	"github.com/begmaroman/acme-dns-route53/utils/strsl"
)

const (
	recordType       = "TXT"
	recordTTL  int64 = 5
)

// Ensures that DNSProvider implements challenge.Provider interface
var _ challenge.Provider = &DNSProvider{}

// DNSProvider is the custom implementation of the challenge.Provider interface
type DNSProvider struct {
	r53 *route53.Route53
}

// NewDNSProviderManual is the constructor of DNSProvider
func NewProvider(r53 *route53.Route53) *DNSProvider {
	return &DNSProvider{
		r53: r53,
	}
}

// Present prints instructions for manually creating the TXT record
func (p *DNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(fqdn)
	if err != nil {
		return errors.Wrapf(err, "unable to find zone by FQDN = '%s'", fqdn)
	}

	log.Infof("[%s] acme: Creating TXT record in %s zone", domain, authZone)

	// Create a subdomain
	recordID, err := p.changeDNSRecord(route53.ChangeActionUpsert, fqdn, p.buildQuotedValue(value))
	if err != nil {
		return errors.Wrapf(err, "unable to change a record with FQDN = '%s'", fqdn)
	}

	log.Infof("[%s] acme: Created TXT record in %s zone with ID %s", domain, authZone, recordID)

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

	log.Infof("[%s] acme: Removing TXT record from %s zone", domain, authZone)

	// Delete the subdomain
	recordID, err := p.changeDNSRecord(route53.ChangeActionDelete, fqdn, p.buildQuotedValue(value))
	if err != nil {
		return errors.Wrapf(err, "unable to delete a record with FQDN = '%s'", fqdn)
	}

	log.Infof("[%s] acme: Removed TXT record in %s zone with ID %s", domain, authZone, recordID)

	return nil
}

// changeDNSRecord changed the record in DNS Route53 by the given params
func (p *DNSProvider) changeDNSRecord(action, domainName, value string) (string, error) {
	tp := recordType
	ttl := recordTTL

	hostedZoneID, err := p.retrieveHostedZone(domainName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve hosted zone ID for domain = '%s'", domainName)
	}

	log.Println("value", value)

	comment := p.buildDNSComment(action, domainName)

	result, err := p.r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &hostedZoneID,
		ChangeBatch: &route53.ChangeBatch{
			Comment: &comment,
			Changes: []*route53.Change{
				{
					Action: &action,
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: &domainName,
						Type: &tp,
						TTL:  &ttl,
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: &value, // TXT record
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return "", errors.Wrapf(err, "unable to change DNS record with HostedZoneId = '%s', Name = '%s', and Action = '%s'", hostedZoneID, domainName, action)
	}

	return *result.ChangeInfo.Id, nil
}

// buildDNSComment creates a comment for changing DNS record
func (p *DNSProvider) buildDNSComment(action, domainName string) string {
	return fmt.Sprintf("acme-dns-route53 certificate validation, action = %s and domain = %s", action, domainName)
}

// buildQuotedValue quotes the given value
func (p *DNSProvider) buildQuotedValue(value string) string {
	return fmt.Sprintf("\"%s\"", value)
}

// retrieveHostedZone retrieves the zone id responsible a given FQDN.
// That is, the id for the zone whose name is the longest parent of the domain.
func (p *DNSProvider) retrieveHostedZone(domainName string) (string, error) {
	zonesList, err := p.r53.ListHostedZones(nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to list hosted zones")
	}

	var zones Zones
	targetLabels := strings.Split(domainName, ".")

	for _, zone := range zonesList.HostedZones {
		// We canno work with a private zone
		if *zone.Config.PrivateZone {
			continue
		}

		if zone.Name == nil {
			continue
		}

		candidateLabels := strings.Split(*zone.Name, ".")
		if strsl.Equal(candidateLabels, targetLabels[len(targetLabels)-len(candidateLabels):]) {
			zones = append(zones, zone)
		}
	}

	if len(zones) == 0 {
		return "", errors.Errorf("unable to find a Route53 hosted zone for domain '%s'", domainName)
	}

	// Sort hosted zones
	sort.Sort(zones)

	return *zones[0].Id, nil
}
