package r53dns

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

// Zones is the custom implementation of sort.Interface for hosted zones list
type Zones []*route53.HostedZone

// Len implements sort.Interface
func (z Zones) Len() int {
	return len(z)
}

// Swap implements sort.Interface
func (z Zones) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}

// Less implements sort.Interface
func (z Zones) Less(i, j int) bool {
	iParts := strings.Split(*z[i].Name, ".")
	jParts := strings.Split(*z[j].Name, ".")

	return len(jParts) < len(iParts)
}
