package transformers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// TransformerWithClosedUpstream decorates a given producer by closing the context before its
// execution.
func TransformerWithClosedUpstream(transformer rivers.Transformer) rivers.Transformer {
	return func(upstream stream.Reader) stream.Reader {
		upstream.Close(nil)
		return transformer(upstream)
	}
}
