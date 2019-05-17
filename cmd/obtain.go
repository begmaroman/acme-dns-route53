package cmd

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/certstore/acmstore"
	"github.com/begmaroman/acme-dns-route53/cmd/flags"
	"github.com/begmaroman/acme-dns-route53/handler"
	"github.com/begmaroman/acme-dns-route53/handler/r53dns"
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

		// Init a common logger
		log := logrus.New()

		// Create a new certificates handler
		h := handler.NewCertificateHandler(&handler.CertificateHandlerOptions{
			ConfigDir:         flags.GetConfigPathFlagValue(cmd),
			Staging:           flags.GetStagingFlagValue(cmd),
			NotificationTopic: flags.GetTopicFlagValue(cmd),
			RenewBefore:       flags.GetRenewBeforeFlagValue(cmd) * 24,
			Log:               log,
			Notifier:          awsns.New(AWSSession, log),    // Initialize SNS API client
			DNS01:             r53dns.New(AWSSession, log),   // Initialize DNS-01 challenge provider by Route 53
			Store:             acmstore.New(AWSSession, log), // Initialize ACM client
		})

		var wg sync.WaitGroup
		for _, domain := range domains {
			wg.Add(1)
			go func(domainList []string) {
				defer wg.Done()

				if err := h.Obtain(domainList, email); err != nil {
					logrus.Errorf("[%s] unable to obtain certificate: %s\n", domain, err)
				}
			}([]string{domain})
		}
		wg.Wait()

		return nil
	},
}

func init() {
	flags.AddDomainsFlag(certificateObtainCmd)
	flags.AddEmailFlag(certificateObtainCmd)
	flags.AddConfigPathFlag(certificateObtainCmd)
	flags.AddStagingFlag(certificateObtainCmd)
	flags.AddTopicFlag(certificateObtainCmd)
	flags.AddRenewBeforeFlag(certificateObtainCmd)

	RootCmd.AddCommand(certificateObtainCmd)
}
