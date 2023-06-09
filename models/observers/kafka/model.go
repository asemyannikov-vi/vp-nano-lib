package observer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	internalchangemanager "github.com/asemyannikov-vi/vp-nano-lib/internals/change_manager"
	internalobserver "github.com/asemyannikov-vi/vp-nano-lib/internals/observer"
)

type observer struct {
	writer *kafka.Writer

	changeManager internalchangemanager.ChangeManager
}

func New(topic string, brokers []string, changeManager interface{}) (internalobserver.Observer, error) {
	castedChangeManager, ok := changeManager.(internalchangemanager.ChangeManager)
	if !ok {
		return nil, errors.New("invalid cast")
	}

	return &observer{
		writer: kafka.NewWriter(
			kafka.WriterConfig{
				Brokers: brokers,
				Topic:   topic,
			},
		),
		changeManager: castedChangeManager,
	}, nil
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
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

	return nil
}
