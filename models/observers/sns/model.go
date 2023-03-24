package observer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	internalchangemanager "vp-nano-lib/internals/change_manager"
	internalobserver "vp-nano-lib/internals/observer"
)

type observer struct {
	sns   *sns.SNS
	topic *string

	changeManager internalchangemanager.ChangeManager
}

func New(topic string, session *session.Session, changeManager interface{}) internalobserver.Observer {
	return &observer{
		sns:   sns.New(session),
		topic: aws.String(topic),

		changeManager: changeManager.(internalchangemanager.ChangeManager),
	}
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
	fmt.Println("Run the SNS Observer")

	valueInBytes, err := observer.changeManager.Manage(context, value)
	if err != nil {
		return err
	}

	if _, err := observer.sns.PublishWithContext(
		*context,
		&sns.PublishInput{
			TopicArn: observer.topic,
			Message:  aws.String(string(valueInBytes)),
		},
	); err != nil {
		return err
	}

	fmt.Println("The event send to the SNS Topic:\n", string(valueInBytes))

	return nil
}
