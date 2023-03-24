package subject

import (
	"context"
	"errors"
	internalobserver "vp-nano-lib/internals/observer"
	internalsubject "vp-nano-lib/internals/subject"
)

type subject struct {
	observers []internalobserver.Observer

	context *context.Context
	value   interface{}
}

func New() internalsubject.Subject {
	return &subject{}
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

func (subject *subject) ListenAndServe() {
}

func (subject *subject) Notify() error {
	for _, observer := range subject.observers {
		value := subject.GetState().(string)
		observer.Update(subject.context, value)
	}

	return nil
}

func (subject *subject) SetState(value interface{}) {
	convertedValue := value.(string)

	subject.value = convertedValue

	subject.Notify()
}

func (subject subject) GetState() interface{} {
	return subject.value
}
