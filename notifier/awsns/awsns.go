package awsns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/begmaroman/acme-dns-route53/notifier"
)

// To make sure that snsNotifier implements notifier.Notifier interface
var _ notifier.Notifier = &snsNotifier{}

// snsNotifier implements notifier.Notifier for ACM by Amazon Web Services
type snsNotifier struct {
	sns *sns.SNS
	log *logrus.Logger
}

// New is the constructor of snsNotifier
func New(provider client.ConfigProvider, log *logrus.Logger) notifier.Notifier {
	return &snsNotifier{
		sns: sns.New(provider),
		log: log,
	}
}

// Notify implements implements notifier.Notifier interface.
// Publishes a message with the given topic to ACM by AWS
func (n *snsNotifier) Notify(topic, message string) error {
	publishResp, err := n.sns.Publish(&sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(message),
	})
	if err != nil {
		return errors.Wrap(err, "sns: unable to publish notification to SNS")
	}

	n.log.Infof("sns: Message with ID '%s' published to topic '%s' successfully", aws.StringValue(publishResp.MessageId), topic)

	return nil
}
