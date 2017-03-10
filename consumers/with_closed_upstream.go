package consumers

import (
	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

// WithClosedUpstream decorates a given producer by closing the context before its
// execution.
func WithClosedUpstream(consumer pipeline.Consumer) pipeline.Consumer {
	return func(upstream stream.Reader) {
		upstream.Close(nil)
		consumer(upstream)
	}
}
