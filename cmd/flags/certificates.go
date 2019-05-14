package flags

import (
	"strings"

	"github.com/spf13/cobra"
)

const (
	domainsSeparator  = ","
	defaultConfigPath = ""
	defaultTopic      = ""

	flagDomains    = "domains"
	flagEmail      = "email"
	flagConfigPath = "config-path"
	flagStaging    = "staging"
	flagTopic      = "topic"
)

// AddDomainsFlag adds the domains flag to the command
func AddDomainsFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagDomains, "", "The domains list, comma-separated", true)
}

// GetDomainsFlagValue gets the value of the domains list from command
func GetDomainsFlagValue(c *cobra.Command) []string {
	domainsString := c.Flag(flagDomains).Value.String()
	return strings.Split(domainsString, domainsSeparator)
}

// AddEmailFlag adds the email flag to the command
func AddEmailFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagEmail, "", "The Email", true)
}

// GetEmailFlagValue gets the value of the email flag from the command
func GetEmailFlagValue(c *cobra.Command) string {
	return c.Flag(flagEmail).Value.String()
}

// AddConfigPathFlag adds the config path flag to the command
func AddConfigPathFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagConfigPath, defaultConfigPath, "The path to config directory", false)
}

// GetConfigPathFlagValue gets the value of the config path flag from the command
func GetConfigPathFlagValue(c *cobra.Command) string {
	return c.Flag(flagConfigPath).Value.String()
}

// AddConfigPathFlag adds the staging flag to the command
func AddStagingFlag(c *cobra.Command) {
	AddPersistentBoolFlag(c, flagStaging, false, "Use --staging flag for using staging Let's Encrypt environment", false)
}

// GetStagingFlagValue gets the value of the staging flag from the command
func GetStagingFlagValue(c *cobra.Command) bool {
	return c.Flag(flagStaging).Value.String() == "true"
}

// AddTopicFlag adds the topic flag to the command
func AddTopicFlag(c *cobra.Command) {
	AddPersistentStringFlag(c, flagTopic, defaultTopic, "Provide SNS notification topic using --topic flag", false)
}

// GetTopicFlagValue gets the value of the topic flag from the command
func GetTopicFlagValue(c *cobra.Command) string {
	return c.Flag(flagTopic).Value.String()
}
