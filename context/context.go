package context

import (
	goContext "context"
)

// Context
//
//  ┏━━━> c2
// c1 ━━> c3 ━━> c5
//  ┗━━━> c4
//
type Context interface {
	goContext.Context
	Open()
	Close()
	Opened() <-chan struct{}
	NewChild() Context
}

func New() Context {
	ctx, cancel := goContext.WithCancel(goContext.Background())
	opened := make(chan struct{})
	return &context{ctx, cancel, opened, make([]Context, 0)}
}

type context struct {
	goContext.Context
	cancel   goContext.CancelFunc
	opened   chan struct{}
	children []Context
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
	child := &context{ctx, cancelWithPropagation, parent.opened, make([]Context, 0)}
	parent.children = append(parent.children, child)
	return child
}
