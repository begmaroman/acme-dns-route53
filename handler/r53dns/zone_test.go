package r53dns

import (
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/stretchr/testify/require"
)

func TestZone(t *testing.T) {
	testTable := []*struct {
		testName      string
		zones         []string
		expectedZones []string
	}{
		{
			testName:      "one subdomain",
			zones:         []string{"test.com", "sub.test.com"},
			expectedZones: []string{"sub.test.com", "test.com"},
		},
		{
			testName:      "two subdomain",
			zones:         []string{"sub1.test.com", "sub2.sub2.test.com", "test.com"},
			expectedZones: []string{"sub2.sub2.test.com", "sub1.test.com", "test.com"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.testName, func(t *testing.T) {
			var zones Zones
			for _, actualZone := range tt.zones {
				zone := actualZone
				zones = append(zones, &route53.HostedZone{
					Name: &zone,
				})
			}

			var expectedZones Zones
			for _, expectedZone := range tt.expectedZones {
				zone := expectedZone
				expectedZones = append(expectedZones, &route53.HostedZone{
					Name: &zone,
				})
			}

			sort.Sort(zones)

			require.Equal(t, expectedZones, zones)
		})
	}
}
