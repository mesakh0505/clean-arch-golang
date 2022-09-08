package sns

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sirupsen/logrus"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
)

// NewEbusSubscriber create the EbusSubscriber for SNS
func NewEbusSubscriber(snsClient *sns.SNS) ebus.HandlerFunc {
	return func(ctx context.Context, e ebus.Event) {
		entry := logrus.WithFields(logrus.Fields{
			"eventName":   e.Name,
			"data":        e.Data,
			"occuredTime": e.OccuredTime,
		})

		// TODO(LieAlbertTriAdrian): ensure the publish input fields here
		// TODO(LieAlbertTriAdrian): Handle Event filtering, to only publish event that should be published to this handler
		msg, err := e.JSONString()
		if err != nil {
			entry.Error("Failed to convert the Event to string")
		}

		input := &sns.PublishInput{
			Message: &msg,
			// TopicArn: , // TODO: (LieAlbertTriAdrian) set this topics
			// TargetArn: ,// TODO: (LieAlbertTriAdrian) set this target
		}

		output, err := snsClient.Publish(input)
		if err != nil {
			entry.Errorf("Failed to Publish event, Got err: %v", err)
		}
		entry.Infof("Event Published Successfully, with ID: %s", *output.MessageId)
	}
}
