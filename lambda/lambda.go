package lambda

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/go-acme/lego/log"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore/acmstore"
	"github.com/begmaroman/acme-dns-route53/handler"
)

const (
	// ConfigDir is the default configuration directory
	ConfigDir = "/tmp"

	// DomainsSeparator is the default domains separator
	DomainsSeparator = ","
)

// AWSSession is the AWS session data
var AWSSession *session.Session

var (
	// ErrEmailMissing is the error when email is not provided
	ErrEmailMissing = errors.New("email must be filled")

	// ErrDomainsMissing is the error when the domains list is empty
	ErrDomainsMissing = errors.New("domains list must not be filled")
)

func init() {
	// Initialize logger
	log.Logger = logrus.New()

	// Initialized AWS session
	AWSSession = session.Must(session.NewSession())
}

func Init() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent() error {
	// Domains list must not be empty
	domains := strings.Split(os.Getenv("LETSENCRYPT_DOMAINS"), DomainsSeparator)
	if len(domains) == 0 {
		return ErrDomainsMissing
	}

	// Email must be filled
	email := os.Getenv("LETSENCRYPT_EMAIL")
	if len(email) == 0 {
		return ErrEmailMissing
	}

	// Check environment
	var isStaging bool
	if os.Getenv("LETSENCRYPT_STAGING") == "1" {
		isStaging = true
	}

	certificateHandler := handler.NewCertificateHandler(
		isStaging,
		acmstore.New(AWSSession),
		route53.New(AWSSession),
		ConfigDir,
	)

	for _, domain := range domains {
		if err := certificateHandler.Obtain([]string{domain}, email); err != nil {
			logrus.Errorf("[%s] unable to obtain certificate: %s", domain, err)
			continue
		}

		logrus.Infof("[%s] certificate successfully obtained and stored", domain)
	}

	return nil
}
