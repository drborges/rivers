package producers

import (
	"github.com/drborges/rivers/context"

	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

// WithClosedContext decorates a given producer by closing the context before its
// execution.
func WithClosedContext(producer pipeline.Producer) pipeline.Producer {
	return func(ctx context.Context) stream.Reader {
		ctx.Close(nil)
		return producer(ctx)
	}
}
