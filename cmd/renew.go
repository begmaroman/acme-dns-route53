package cmd

import (
	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/cmd/flags"
	"github.com/begmaroman/acme-dns-route53/handler"
)

// certificateRenewCmd represents the certificate renewing command
var certificateRenewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew SSL certificates",
	Long:  `This command renews SSL certificates by the given domains using the given email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Inits arguments
		domains := flags.GetDomainsFlagValue(cmd)
		email := flags.GetEmailFlagValue(cmd)

		// Create a new certificates handler
		h := handler.NewCertificateHandler(ResultPrinter, FileStore, Route53)

		return h.Renew(domains, email)
	},
}

func init() {
	flags.AddDomainsFlag(certificateRenewCmd)
	flags.AddEmailFlag(certificateRenewCmd)

	RootCmd.AddCommand(certificateRenewCmd)
}
