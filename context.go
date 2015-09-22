package rivers

import (
	"errors"
	"fmt"
	"github.com/drborges/rivers/stream"
	"runtime/debug"
	"time"
)

var DebugEnabled = false

type context struct {
	requests chan error
	success  chan struct{}
	failure  chan struct{}
	deadline time.Duration
	err      error
}

func NewContext() stream.Context {
	return &context{
		requests: make(chan error, 1),
		success:  make(chan struct{}),
		failure:  make(chan struct{}),
		deadline: time.Hour,
	}
}

func (context *context) Err() error {
	return context.err
}

func (context *context) Deadline() time.Duration {
	return context.deadline
}

func (context *context) SetDeadline(duration time.Duration) {
	context.deadline = duration
}

func (context *context) Failure() <-chan struct{} {
	return context.failure
}

func (context *context) Done() <-chan struct{} {
	return context.success
}

func (context *context) Close(err error) {
	context.requests <- err
	ch := context.success
	if err != nil {
		ch = context.failure
	}

	select {
	case <-ch:
		return
	default:
		close(ch)
		context.err = <-context.requests
	}
}

func (context *context) Recover() {
	if r := recover(); r != nil {
		if DebugEnabled && r != stream.Done {
			debug.PrintStack()
		}
		err := errors.New(fmt.Sprintf("Recovered from %v", r))
		if e, ok := r.(error); ok {
			err = e
		}
		context.Close(err)
	}
}
