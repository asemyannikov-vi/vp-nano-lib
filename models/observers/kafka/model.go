package observer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	internalchangemanager "vp-nano-lib/internals/change_manager"
	internalobserver "vp-nano-lib/internals/observer"
)

type observer struct {
	writer *kafka.Writer

	changeManager internalchangemanager.ChangeManager
}

func New(topic string, brokers []string, changeManager interface{}) internalobserver.Observer {
	return &observer{
		writer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers: brokers,
				Topic:   topic,
			},
		),

		changeManager: changeManager.(internalchangemanager.ChangeManager),
	}
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
	fmt.Println("Run the Kafka Observer")

	valueInBytes, err := observer.changeManager.Manage(context, value)
	if err != nil {
		return err
	}

	if err := observer.writer.WriteMessages(
		*context,
		kafka.Message{
			Key:   []byte(uuid.New().String()),
			Value: []byte(valueInBytes),
		},
	); err != nil {
		return err
	}

	fmt.Println("The event sent to the Kafka Topic:\n", string(valueInBytes))

	return nil
}
