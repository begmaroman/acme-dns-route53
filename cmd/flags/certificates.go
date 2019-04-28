package flags

import (
	"strings"

	"github.com/spf13/cobra"
)

const (
	domainsSeparator = ","

	flagDomains = "domains"
	flagEmail   = "email"
)

// AddDomainsFlag adds the domains flag to the command
func AddDomainsFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagDomains, "", "The domains list, comma-separated", true)
}

// GetDomainsFlagValue gets the domains list from command
func GetDomainsFlagValue(c *cobra.Command) []string {
	domainsString := c.Flag(flagDomains).Value.String()
	return strings.Split(domainsString, domainsSeparator)
}

// AddEmailFlag adds the email flag to the command
func AddEmailFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagEmail, "", "The Email", true)
}

// GetEmailFlagValue gets the email flag from the command
func GetEmailFlagValue(c *cobra.Command) string {
	return c.Flag(flagEmail).Value.String()
}
