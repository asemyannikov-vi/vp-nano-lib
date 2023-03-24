package observer

import (
	"context"
	"fmt"

	internalobserver "vp-nano-lib/internals/observer"
)

type observer struct {
	Name string
}

func New(name string) internalobserver.Observer {
	return &observer{
		Name: name,
	}
}

func (observer *observer) Update(context *context.Context, value interface{}) error {
	concreteValue := value.(string)
	fmt.Println("Observer", observer.Name, "got value", concreteValue)
	return nil
}
