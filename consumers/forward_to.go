package consumers

import (
	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

// ForwardTo implements a pipeline.Consumer that forwards stream data to a given
// writable stream.
func ForwardTo(ch stream.Writable) pipeline.Consumer {
	return func(upstream stream.Reader) {
		go func() {
			defer upstream.Close(nil)
			defer close(ch)

			for data := range upstream.Read() {
				ch <- data
			}
		}()
	}
}
