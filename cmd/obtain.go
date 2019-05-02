package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/cmd/flags"
	"github.com/begmaroman/acme-dns-route53/handler"
)

// certificateObtainCmd represents the certificate obtaining command
var certificateObtainCmd = &cobra.Command{
	Use:   "obtain",
	Short: "Obtain SSL certificates",
	Long:  `This command creates new SSL certificates or renews existing ones for the given domains using the given parameters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Inits arguments
		isStaging := flags.GetStagingFlagValue(cmd)
		domains := flags.GetDomainsFlagValue(cmd)
		email := flags.GetEmailFlagValue(cmd)
		configPath := flags.GetConfigPathFlagValue(cmd)

		// Create a new certificates handler
		h := handler.NewCertificateHandler(isStaging, CertStore, Route53, configPath)

		for _, domain := range domains {
			if err := h.Obtain([]string{domain}, email); err != nil {
				logrus.Errorf("[%s] unable to obtain certificate: %s\n", domain, err)
				continue
			}

			logrus.Infof("[%s] certificate successfully obtained and stored\n", domain)
		}

		return nil
	},
}

func init() {
	flags.AddDomainsFlag(certificateObtainCmd)
	flags.AddEmailFlag(certificateObtainCmd)
	flags.AddConfigPathFlag(certificateObtainCmd)
	flags.AddStagingFlag(certificateObtainCmd)

	RootCmd.AddCommand(certificateObtainCmd)
}
