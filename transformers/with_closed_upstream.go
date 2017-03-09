package transformers

import (
	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

// WithClosedUpstream decorates a given producer by closing the context before its
// execution.
func WithClosedUpstream(transformer pipeline.Transformer) pipeline.Transformer {
	return func(upstream stream.Reader) stream.Reader {
		upstream.Close(nil)
		return transformer(upstream)
	}
}
