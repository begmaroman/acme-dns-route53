package r53dns

import (
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/begmaroman/acme-dns-route53/utils/strsl"
	"github.com/go-acme/lego/log"
	"github.com/pkg/errors"
)

const (
	// recordTTL is the default TTL (in seconds)
	recordTTL int64 = 60
)

// r53ResourceWorker represents the functionality to work with Route53 API
type r53ResourceWorker struct {
	r53 *route53.Route53
}

// newR53ResourceWorker is the constructor of r53ResourceWorker
func newR53ResourceWorker(r53 *route53.Route53) *r53ResourceWorker {
	return &r53ResourceWorker{
		r53: r53,
	}
}

// changeDNSRecord changed the record in DNS Route53 by the given params
func (r *r53ResourceWorker) changeDNSRecord(action, domainName, value string) (string, error) {
	// Retrieve a hosted zone ID
	hostedZoneID, err := r.retrieveHostedZone(domainName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve hosted zone ID for domain = '%s'", domainName)
	}

	log.Infof("[%s] acme: Changing record (action '%s') in the zone with ID = '%s'", domainName, action, hostedZoneID)

	// Build comment for the current action
	comment := buildDNSComment(action, domainName)

	// Change the record
	result, err := r.r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String(comment),
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domainName),
						Type: aws.String(route53.RRTypeTxt),
						TTL:  aws.Int64(recordTTL),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(value), // TXT record
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

	// Wait for change
	if err := r.waitForChange(*result.ChangeInfo.Id); err != nil {
		return "", errors.Wrap(err, "failed while waiting for change status")
	}

	return *result.ChangeInfo.Id, nil
}

// retrieveHostedZone retrieves the zone id responsible a given FQDN.
// That is, the id for the zone whose name is the longest parent of the domain.
func (r *r53ResourceWorker) retrieveHostedZone(domainName string) (string, error) {
	zonesList, err := r.r53.ListHostedZones(nil)
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

// waitForChange waits for a change to be propagated to all Route53 DNS servers
func (r *r53ResourceWorker) waitForChange(changeID string) error {
	// Check change
	for i := 0; i < 120; i++ {
		changeResp, err := r.r53.GetChange(&route53.GetChangeInput{
			Id: aws.String(changeID),
		})
		if err != nil {
			return errors.Wrap(err, "unable to get changing status")
		}

		if *changeResp.ChangeInfo.Status == route53.ChangeStatusInsync {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}
