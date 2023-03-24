package changemanager

import "context"

type ChangeManager interface {
	Manage(context *context.Context, value interface{}) ([]byte, error)
}
