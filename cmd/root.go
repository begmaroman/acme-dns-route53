package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/spf13/cobra"

	"github.com/begmaroman/acme-dns-route53/certstore"
	"github.com/begmaroman/acme-dns-route53/printer"
)

var (
	// ResultPrinter is the printer used to print command results and errors
	ResultPrinter printer.Printer = printer.NewStandardOutputPrinter(os.Stdout)

	// FileStore is the store used for CRUD operations with certificates
	FileStore certstore.CertStore = certstore.NewFileStore()

	// Route53 is the Route53 client from just a session.
	// Initial credentials loaded from SDK's default credential chain. Such as
	// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
	// Role. These credentials will be used to to make the STS Assume Role API.
	Route53 *route53.Route53 = route53.New(session.Must(session.NewSession()))

	// RootCmd represents the base command when called without any subcommands
	RootCmd = &cobra.Command{
		Use:   "acme-dns-route53",
		Short: "DNS-01 challenge resolver",
		Long:  `acme-dns-route53 is a CLI for managing SSL certificates using DNS-01 challenge.`,
	}
)

func init() {
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
