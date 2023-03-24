package observer

import (
	"context"
	"errors"

	internalchangemanager "github.com/asemyannikov-vi/vp-nano-lib/internals/change_manager"
	internalobserver "github.com/asemyannikov-vi/vp-nano-lib/internals/observer"
)

type observer struct {
	changeManager internalchangemanager.ChangeManager
}

func New(changeManager interface{}) (internalobserver.Observer, error) {
	castedChangeManager, ok := changeManager.(internalchangemanager.ChangeManager)
	if !ok {
		return nil, errors.New("invalid cast")
	}

	return &observer{
		changeManager: castedChangeManager,
	}, nil
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
	_, err := observer.changeManager.Manage(context, value)
	if err != nil {
		return err
	}

	return nil
}
