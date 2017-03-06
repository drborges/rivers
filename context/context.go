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
	// Open opens the context and all its descendants so that taks bound to these
	// contexts can start their work.
	Open()
	// Close attempts to close the context. If the context still has opened
	// children, this operation will be a noop.
	Close()
	// Opened read-only channel used to check whether or not the context is still
	// opened for work.
	Opened() <-chan struct{}
	// NewChild creates a new child Context.
	NewChild() Context
}

// New creates a new Context.
func New() Context {
	ctx, cancel := goContext.WithCancel(goContext.Background())
	opened := make(chan struct{})
	return &context{ctx, cancel, opened, make([]Context, 0), DefaultConfig}
}

type context struct {
	goContext.Context
	cancel   goContext.CancelFunc
	opened   chan struct{}
	children []Context
	config   Config
}

func (ctx *context) Config() Config {
	return ctx.config
}

func (ctx *context) Close() {
	for _, child := range ctx.children {
		select {
		case <-child.Done():
		default:
			// Parent can only close the context when all children have closed theirs.
			return
		}
	}
	ctx.cancel()
}

func (ctx *context) Open() {
	close(ctx.opened)
}

func (ctx *context) Opened() <-chan struct{} {
	return ctx.opened
}

func (parent *context) NewChild() Context {
	ctx, cancel := goContext.WithCancel(parent.Context)
	cancelWithPropagation := func() {
		cancel()
		parent.Close()
	}
	child := &context{ctx, cancelWithPropagation, parent.opened, make([]Context, 0), parent.config}
	parent.children = append(parent.children, child)
	return child
}
