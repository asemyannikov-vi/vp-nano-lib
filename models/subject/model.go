package subject

import (
	"context"
	"errors"
	"fmt"

	internalobserver "github.com/asemyannikov-vi/vp-nano-lib/internals/observer"
	internalsubject "github.com/asemyannikov-vi/vp-nano-lib/internals/subject"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

const (
	SUBJECT_KAFKA   = "kafka"
	SUBJECT_AWS_SQS = "awssqs"
)

type KafkaTopicSettings struct {
	Topic   string
	Address string
}

type AWSSQSTopicSettings struct {
	Topic   string
	Region  string
	Account string
}

type subject struct {
	observers []internalobserver.Observer

	subscription *pubsub.Subscription
	context      context.Context
	writeChannel chan []byte
	readChannel  chan []byte
}

func New(
	context context.Context,
	driver string,
	topicSettings interface{},
) (internalsubject.Subject, error) {
	var url string

	if driver == SUBJECT_KAFKA {
		kafkaTopicSettings := topicSettings.(KafkaTopicSettings)
		url = fmt.Sprintf("kafka://%s?topic=%s", kafkaTopicSettings.Address, kafkaTopicSettings.Topic)
	} else if driver == SUBJECT_AWS_SQS {
		awssqsTopicSettings := topicSettings.(AWSSQSTopicSettings)
		url = fmt.Sprintf("awssqs://sqs.%s.amazonaws.com/%s/%s", awssqsTopicSettings.Region, awssqsTopicSettings.Account, awssqsTopicSettings.Topic)
	} else {
		return nil, errors.New("unknown subject driver")
	}

	subscription, err := pubsub.OpenSubscription(context, url)
	if err != nil {
		return nil, err
	}

	return &subject{
		context:      context,
		subscription: subscription,
		writeChannel: make(chan []byte),
		readChannel:  make(chan []byte),
	}, nil
}

func (subject *subject) Attach(newObserver internalobserver.Observer) error {
	for _, observer := range subject.observers {
		if observer == newObserver {
			return errors.New("observer already attached")
		}
	}

	subject.observers = append(subject.observers, newObserver)
	return nil
}

func (subject *subject) Detach(newObserver internalobserver.Observer) error {
	for i, observer := range subject.observers {
		if observer == newObserver {
			subject.observers = append(subject.observers[:i], subject.observers[i+1:]...)
			return nil
		}
	}

	return errors.New("observer not found")
}

func (subject *subject) Notify(value []byte) error {
	for _, observer := range subject.observers {
		observer.Update(&subject.context, value)
	}

	return nil
}

func (subject *subject) poll() {
	for {
		message, err := subject.subscription.Receive(subject.context)
		if err != nil {
			continue
		}

		subject.Notify(message.Body)

		message.Ack()

		subject.SetChannelState(message.Body)
	}
}

func (subject *subject) ListenAndServe() {
	go subject.poll()
}

func (subject subject) SetChannelState(value []byte) {
	subject.writeChannel <- value
}

func (subject subject) GetChannelState() []byte {
	return <-subject.readChannel
}

func (subject subject) Monitor() {
	for {
		writeValue := <-subject.writeChannel
		value := writeValue
		subject.readChannel <- value
	}
}
