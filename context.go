package rivers

import (
	"errors"
	"fmt"
	"github.com/drborges/rivers/stream"
	"runtime/debug"
)

var DebugEnabled = false

type context struct {
	requests chan error
	closed   chan struct{}
	err      error
}

func NewContext() stream.Context {
	return &context{
		requests: make(chan error, 1),
		closed:   make(chan struct{}),
	}
}

func (context *context) Close(err error) {
	context.requests <- err
	select {
	// TODO timeout?
	case <-context.closed:
		return
	default:
		close(context.closed)
		context.err = <-context.requests
	}
}

func (context *context) Err() error {
	return context.err
}

func (context *context) Failure() <-chan struct{} {
	return context.closed
}

func (context *context) Recover() {
	if r := recover(); r != nil && r != stream.Done {
		if DebugEnabled {
			debug.PrintStack()
		}
		err := errors.New(fmt.Sprintf("Recovered from %v", r))
		if e, ok := r.(error); ok {
			err = e
		}
		context.Close(err)
	}
}
