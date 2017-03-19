package aggregators

import (
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/stream"
)

// FIFO merges two streams into one applying a first-in-first-out strategy.
// Implements rivers.Aggregator interface
func FIFO(upstream1, upstream2 stream.Reader) stream.Reader {
	reader, writer := stream.New(ctxtree.New())

	upstream1Done := make(chan struct{}, 0)
	upstream2Done := make(chan struct{}, 0)

	// Forwards data from upstream1 to downstream.
	go func() {
		defer close(upstream1Done)

		for data := range upstream1.Read() {
			if err := writer.Write(data); err != nil {
				return
			}
		}
	}()

	// Forwards data from upstream2 to downstream.
	go func() {
		defer close(upstream2Done)

		for data := range upstream2.Read() {
			if err := writer.Write(data); err != nil {
				return
			}
		}
	}()

	// Propagates downstream cancellation to the upstream.
	go func() {
		defer upstream1.Close(nil)
		defer upstream2.Close(nil)
		<-reader.Done()
	}()

	// Closes downstream when there is no longer upstream data to be consumed.
	go func() {
		defer writer.Close(nil)
		<-upstream1Done
		<-upstream2Done
	}()

	return reader
}
