package flags

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AddEnvVarPersistentFlag adds a flag that can also be passed via environment variable to the command
func AddEnvVarPersistentFlag(c *cobra.Command, flag string, envVar string, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().String(flag, os.Getenv(envVar), fmt.Sprintf("%s [env=%s]%s", description, envVar, req))

	if isRequired && os.Getenv(envVar) == "" {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentStringFlag adds a string flag to the command
func AddPersistentStringFlag(c *cobra.Command, flag string, value string, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().String(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentBoolFlag adds a bool flag to the command
func AddPersistentBoolFlag(c *cobra.Command, flag string, value bool, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().Bool(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}

// AddPersistentIntFlag adds a int flag to the command
func AddPersistentIntFlag(c *cobra.Command, flag string, value int, description string, isRequired bool) {
	req := ""
	if isRequired {
		req = " (required)"
	}

	c.PersistentFlags().Int(flag, value, fmt.Sprintf("%s%s", description, req))

	if isRequired {
		c.MarkPersistentFlagRequired(flag)
	}
}
