package lambda

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/certstore/acmstore"
	"github.com/begmaroman/acme-dns-route53/handler"
	"github.com/begmaroman/acme-dns-route53/handler/r53dns"
	"github.com/begmaroman/acme-dns-route53/notifier/awsns"
)

const (
	// ConfigDir is the default configuration directory
	ConfigDir = "/tmp"
)

var (
	// ErrEmailMissing is the error when email is not provided
	ErrEmailMissing = errors.New("email must be filled")

	// ErrDomainsMissing is the error when the domains list is empty
	ErrDomainsMissing = errors.New("domains list must not be filled")
)

// Payload contains payload data
type Payload struct {
	Domains     []string `json:"domains"`
	Email       string   `json:"email"`
	Staging     string   `json:"staging"`
	Topic       string   `json:"topic"`
	RenewBefore int      `json:"renew_before"`
}

func HandleLambdaEvent(payload Payload) error {
	conf := InitConfig(payload)

	// Domains list must not be empty
	if len(conf.Domains) == 0 {
		return ErrDomainsMissing
	}

	// Email must be filled
	if len(conf.Email) == 0 {
		return ErrEmailMissing
	}

	log := logrus.New()

	// Create a new handler
	certificateHandler := handler.NewCertificateHandler(&handler.CertificateHandlerOptions{
		ConfigDir:         ConfigDir,
		Staging:           conf.Staging,
		NotificationTopic: conf.Topic,
		RenewBefore:       conf.RenewBefore * 24,
		Log:               log,
		Notifier:          awsns.New(AWSSession, log),    // Initialize SNS API client
		DNS01:             r53dns.New(AWSSession, log),   // Initialize DNS-01 challenge provider by Route 53
		Store:             acmstore.New(AWSSession, log), // Initialize ACM client
	})

	var wg sync.WaitGroup
	for _, domain := range conf.Domains {
		wg.Add(1)
		go func(domainList []string) {
			defer wg.Done()

			if err := certificateHandler.Obtain(domainList, conf.Email); err != nil {
				logrus.Errorf("[%s] unable to obtain certificate: %s\n", domain, err)
			}
		}([]string{domain})
	}
	wg.Wait()

	return nil
}
