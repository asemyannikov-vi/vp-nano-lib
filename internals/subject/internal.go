package subject

import (
	"vp-nano-lib/internals/observer"
)

type Subject interface {
	Attach(observer observer.Observer) error
	Detach(observer observer.Observer) error
	Notify() error
	SetState(value interface{})
	GetState() interface{}
	ListenAndServe()
}
