package consumers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// ConsumerWithClosedUpstream decorates a given producer by closing the context before its
// execution.
func ConsumerWithClosedUpstream(consumer rivers.Consumer) rivers.Consumer {
	return func(upstream stream.Reader) {
		upstream.Close(nil)
		consumer(upstream)
	}
}
