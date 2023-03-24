package observer

import (
	"context"

	internalchangemanager "vp-nano-lib/internals/change_manager"
	internalobserver "vp-nano-lib/internals/observer"
)

type observer struct {
	changeManager internalchangemanager.ChangeManager
}

func New(changeManager interface{}) internalobserver.Observer {
	return &observer{
		changeManager: changeManager.(internalchangemanager.ChangeManager),
	}
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
	_, err := observer.changeManager.Manage(context, value)
	if err != nil {
		return err
	}

	return nil
}
