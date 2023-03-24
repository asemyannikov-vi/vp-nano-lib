package sqs

import (
	"context"
	"errors"
	"time"
	internalobserver "vp-nano-lib/internals/observer"
	internalsubject "vp-nano-lib/internals/subject"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

type subject struct {
	observers []internalobserver.Observer

	message      []byte
	pollInterval int64
	subscription *pubsub.Subscription
	context      context.Context
}

func New(
	context context.Context,
	pollInterval int64,
	address string,
) (internalsubject.Subject, error) {
	subscription, err := pubsub.OpenSubscription(context, address)
	if err != nil {
		return nil, err
	}

	return &subject{
		context:      context,
		subscription: subscription,
		pollInterval: pollInterval,
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

func (subject *subject) poll() {
	for {
		time.Sleep(time.Second * time.Duration(subject.pollInterval))

		message, err := subject.subscription.Receive(subject.context)
		if err != nil {
			continue
		}

		subject.SetState(message.Body)

		message.Ack()
	}
}

func (subject *subject) ListenAndServe() {
	go subject.poll()
}

func (subject *subject) Notify() error {
	for _, observer := range subject.observers {
		observer.Update(&subject.context, subject.GetState())
	}

	return nil
}

func (subject *subject) SetState(value interface{}) {
	subject.message = value.([]byte)

	subject.Notify()
}

func (subject subject) GetState() interface{} {
	return subject.message
}
