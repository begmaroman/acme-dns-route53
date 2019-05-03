package lambda

import (
	"errors"

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
)

// AWSSession is the AWS session data
var AWSSession *session.Session

var (
	// ErrEmailMissing is the error when email is not provided
	ErrEmailMissing = errors.New("email must be filled")

	// ErrDomainsMissing is the error when the domains list is empty
	ErrDomainsMissing = errors.New("domains list must not be filled")
)

// Params contains configuration data
type Params struct {
	Domains []string `json:"domains"`
	Email   string   `json:"email"`
	Staging bool     `json:"staging"`
}

func init() {
	// Initialize logger
	log.Logger = logrus.New()

	// Initialized AWS session
	AWSSession = session.Must(session.NewSession())
}

func Init() {
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(params Params) error {
	// Domains list must not be empty
	if len(params.Domains) == 0 {
		return ErrDomainsMissing
	}

	// Email must be filled
	if len(params.Email) == 0 {
		return ErrEmailMissing
	}

	certificateHandler := handler.NewCertificateHandler(
		params.Staging,
		acmstore.New(AWSSession),
		route53.New(AWSSession),
		ConfigDir,
	)

	for _, domain := range params.Domains {
		if err := certificateHandler.Obtain([]string{domain}, params.Email); err != nil {
			logrus.Errorf("[%s] unable to obtain certificate: %s", domain, err)
			continue
		}

		logrus.Infof("[%s] certificate successfully obtained and stored", domain)
	}

	return nil
}
