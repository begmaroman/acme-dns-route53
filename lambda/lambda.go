package lambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-acme/lego/log"
	"github.com/sirupsen/logrus"
)

// AWSSession is the AWS session data
var AWSSession *session.Session

func init() {
	// Initialize logger
	log.Logger = logrus.New()

	// Initialized AWS session
	AWSSession = session.Must(session.NewSession())
}

func Init() {
	lambda.Start(HandleLambdaEvent)
}
