package splitters

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// SplitterWithClosedUpstream decorates a given producer by closing the context before its
// execution.
func SplitterWithClosedUpstream(splitter rivers.Splitter) rivers.Splitter {
	return func(upstream stream.Reader) (stream.Reader, stream.Reader) {
		upstream.Close(nil)
		return splitter(upstream)
	}
}
