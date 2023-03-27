package subject

import (
	"github.com/asemyannikov-vi/vp-nano-lib/internals/observer"
)

type Subject interface {
	Attach(observer observer.Observer) error
	Detach(observer observer.Observer) error
	Notify(value []byte) error

	ListenAndServe()
	Monitor()
	SetChannelState(value []byte)
	GetChannelState() []byte
}
