package r53dns

import "fmt"

// buildQuotedValue quotes the given value
func buildQuotedValue(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

// buildDNSComment creates a comment for changing DNS record
func buildDNSComment(action, domainName string) string {
	return fmt.Sprintf("acme-dns-route53 certificate validation, action = %s and domain = %s", action, domainName)
}
