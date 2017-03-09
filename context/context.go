package context

import (
	goContext "context"
	"time"
)

var (
	// DefaultConfig the default configuration used to create a context.
	DefaultConfig = Config{
		Timeout:    5 * time.Second,
		BufferSize: 1000,
	}
)

type closeFunc func(error)

// Config configuration used to adjust the behavior of the context.
type Config struct {
	// Timeout the timeout which after the context is automatically closed.
	Timeout time.Duration
	// BufferSize the size of the buffer used by workers bound to this context.
	BufferSize int
}

// Context implements the golang context.Context interface, with support to
// Context Trees, providing a different semanthics for cancellation propagation.
// In a Context Tree, cancellation is propagated from leaf nodes up to their
// parent. A parent node can only be canceled if all its children are canceled.
// For example, given the following Context Tree:
//
//  ┏━━━> c2
// c1 ━━> c3 ━━> c5
//  ┗━━━> c4
//
// The context c1, con only be canceld if all its children (c2, c3, and c4) are
// already canceld. Similarly, c3 can only be canceled if c5 is canceled,
// threfore, c1 depends on c5 being canceled befre it may be canceled.
//
// This abstraction enables the creation of stream pipelines, wheren downstreams
// can signal to their upstream when they are done consuming data, freeing the
// upstream to cease its work when no more data is required by its downstreams.
type Context interface {
	// Implements the golang context.Context interface.
	goContext.Context
	// Config returns the configuration used to create the context.
	Config() Config
	// Close attempts to close the context. If the context still has opened
	// children, this operation will be a noop, unless an error is provided, which
	// then causes the whole context tree to be closed.
	Close(error)
	// NewChild creates a new child Context.
	NewChild() Context
}

// New creates a new Context.
func New() Context {
	return FromStdContext(goContext.Background())
}

// FromStdContext creates a new Context from the standard golang context.
func FromStdContext(stdCtx goContext.Context) Context {
	ctx, cancel := goContext.WithCancel(stdCtx)
	context := &context{
		Context:  ctx,
		config:   DefaultConfig,
		children: make([]Context, 0),
	}

	closeFunc := func(err error) {
		cancel()
		context.err = err
	}

	context.closeFunc = closeFunc
	return context
}

// WithConfig creates a new Context with the provided configuration.
func WithConfig(parent Context, config Config) Context {
	ctx, cancel := goContext.WithTimeout(parent, config.Timeout)
	closeFunc := func(err error) {
		cancel()
		parent.Close(err)
	}

	return &context{
		Context:   ctx,
		config:    config,
		closeFunc: closeFunc,
		children:  make([]Context, 0),
	}
}

type context struct {
	goContext.Context
	closeFunc closeFunc
	children  []Context
	config    Config
	err       error
}

func (ctx *context) Config() Config {
	return ctx.config
}

func (ctx *context) Close(err error) {
	if err == nil {
		for _, child := range ctx.children {
			select {
			case <-child.Done():
			default:
				// Parent can only close the context when all children have closed theirs.
				return
			}
		}
	}

	ctx.closeFunc(err)
}

func (ctx *context) NewChild() Context {
	childCtx, cancel := goContext.WithCancel(ctx.Context)
	closeFunc := func(err error) {
		cancel()
		ctx.Close(err)
	}

	child := &context{
		Context:   childCtx,
		closeFunc: closeFunc,
		config:    ctx.config,
		children:  make([]Context, 0),
	}

	ctx.children = append(ctx.children, child)
	return child
}

func (ctx *context) Err() error {
	if ctx.err != nil {
		return ctx.err
	}

	return ctx.Context.Err()
}
