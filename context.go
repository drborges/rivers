package rivers

import (
	"errors"
	"fmt"
	"github.com/drborges/riversv2/rx"
)

type context struct {
	closed chan error
	err    error
}

func NewContext() rx.Context {
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
		c.err = errors.New(fmt.Sprintf("Recovered from %v", r))
		if e, ok := r.(error); ok {
			c.err = e
		}
	}
}
