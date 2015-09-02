package rivers

import (
	"errors"
	"fmt"
	"github.com/drborges/rivers/stream"
	"runtime/debug"
)

var DebugEnabled = false

type context struct {
	closed chan error
	err    error
}

func NewContext() stream.Context {
	return &context{
		closed: make(chan error),
	}
}

func (c *context) Close() {
	select {
	case <-c.closed:
		return
	default:
		close(c.closed)
	}
}

func (c *context) Err() error {
	return c.err
}

func (c *context) Closed() <-chan error {
	return c.closed
}

func (c *context) Recover() {
	if r := recover(); r != nil {
		c.Close()
		if DebugEnabled {
			debug.PrintStack()
		}
		c.err = errors.New(fmt.Sprintf("Recovered from %v", r))
		if e, ok := r.(error); ok {
			c.err = e
		}
	}
}
