package awsns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"

	"github.com/begmaroman/acme-dns-route53/notifier"
)

// To make sure that snsNotifier implements notifier.Notifier interface
var _ notifier.Notifier = &snsNotifier{}

// snsNotifier implements notifier.Notifier for ACM by Amazon Web Services
type snsNotifier struct {
	sns *sns.SNS
}

// New is the constructor of snsNotifier
func New(sns *sns.SNS) notifier.Notifier {
	return &snsNotifier{
		sns: sns,
	}
}

// Notify implements implements notifier.Notifier interface.
// Publishes a message with the given topic to ACM by AWS
func (n *snsNotifier) Notify(topic, message string) error {
	_, err := n.sns.Publish(&sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(message),
	})
	if err != nil {
		return errors.Wrap(err, "unable to publish notification to SNS")
	}

	return nil
}
