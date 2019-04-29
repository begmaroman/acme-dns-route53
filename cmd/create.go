package cmd

import (
	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/cmd/flags"
	"github.com/begmaroman/acme-dns-route53/handler"
)

// certificateCreateCmd represents the certificate creation command
var certificateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create SSL certificates",
	Long:  `This command creates new SSL certificates for the given domains using the given email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Inits arguments
		domains := flags.GetDomainsFlagValue(cmd)
		email := flags.GetEmailFlagValue(cmd)

		// Create a new certificates handler
		h := handler.NewCertificateHandler(ResultPrinter, CertStore, Route53)

		return h.Create(domains, email)
	},
}

func init() {
	flags.AddDomainsFlag(certificateCreateCmd)
	flags.AddEmailFlag(certificateCreateCmd)

	RootCmd.AddCommand(certificateCreateCmd)
}
