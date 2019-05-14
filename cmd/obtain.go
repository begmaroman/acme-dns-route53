package cmd

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/certstore/acmstore"
	"github.com/begmaroman/acme-dns-route53/cmd/flags"
	"github.com/begmaroman/acme-dns-route53/handler"
	"github.com/begmaroman/acme-dns-route53/notifier/awsns"
)

// certificateObtainCmd represents the certificate obtaining command
var certificateObtainCmd = &cobra.Command{
	Use:   "obtain",
	Short: "Obtain SSL certificates",
	Long:  `This command creates new SSL certificates or renews existing ones for the given domains using the given parameters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Inits needed parameters
		domains := flags.GetDomainsFlagValue(cmd)
		email := flags.GetEmailFlagValue(cmd)

		// Create handler options
		certificateHandlerOpts := &handler.CertificateHandlerOptions{
			ConfigDir:         flags.GetConfigPathFlagValue(cmd),
			Staging:           flags.GetStagingFlagValue(cmd),
			NotificationTopic: flags.GetTopicFlagValue(cmd),
			Log:               logrus.New(),                           // Create a new logger
			Notifier:          awsns.New(AWSSession, logrus.New()),    // Initialize SNS API client
			R53:               route53.New(AWSSession),                // Initialize Route53 API client
			Store:             acmstore.New(AWSSession, logrus.New()), // Initialize ACM client
		}

		// Create a new certificates handler
		h := handler.NewCertificateHandler(certificateHandlerOpts)

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
	flags.AddTopicFlag(certificateObtainCmd)

	RootCmd.AddCommand(certificateObtainCmd)
}
