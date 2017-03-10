package producers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/stream"
)

// WithClosedContext decorates a given producer by closing the context before its
// execution.
func WithClosedContext(producer rivers.Producer) rivers.Producer {
	return func(ctx ctxtree.Context) stream.Reader {
		ctx.Close(nil)
		return producer(ctx)
	}
}
