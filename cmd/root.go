package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-acme/lego/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// AWSSession is the session of Amazon Web Services.
	// Initial credentials loaded from SDK's default credential chain. Such as
	// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
	// Role. These credentials will be used to to make the STS Assume Role API.
	AWSSession *session.Session

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "acme-dns-route53",
		Short: "DNS-01 challenge resolver",
		Long:  `acme-dns-route53 is a CLI for managing SSL certificates using DNS-01 challenge.`,
	}
)

func init() {
	// Initialize logger
	log.Logger = logrus.New()

	// Initialized AWS session
	AWSSession = session.Must(session.NewSession())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
