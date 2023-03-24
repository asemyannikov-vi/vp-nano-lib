package observer

import "context"

type Observer interface {
	Update(context *context.Context, value interface{}) error
}
