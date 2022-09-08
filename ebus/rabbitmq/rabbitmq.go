package rabbitmq

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
)

// NewEbusSubscriber create the EbusSubscriber for RabbitMQ
func NewEbusSubscriber(rabbitClient *amqp.Channel) ebus.HandlerFunc {
	return func(ctx context.Context, e ebus.Event) {
		entry := logrus.WithFields(logrus.Fields{
			"eventName":   e.Name,
			"data":        e.Data,
			"occuredTime": e.OccuredTime,
		})

		// TODO(LieAlbertTriAdrian): ensure the publish input fields here
		// TODO(LieAlbertTriAdrian): Handle Event filtering, to only publish event that should be published to this handler
		eventStr, err := e.JSONString()
		if err != nil {
			entry.Error("Failed to convert the Event to string")
		}
		exhcange := ""
		key := ""
		mandatory := false
		immediate := true

		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(eventStr),
			Timestamp:   e.OccuredTime,
			// TODO (LieAlbertTriAdrian): complete the message body
		}
		err = rabbitClient.Publish(exhcange, key, mandatory, immediate, msg)
		if err != nil {
			logrus.Error(err)
		}
		if err != nil {
			entry.Errorf("Failed to Publish event, Got err: %v", err)
		}
		entry.Info("Event Published Successfully to RabbitMq")
	}
}
