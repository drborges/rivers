package consumers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// ForwardTo implements a rivers.Consumer that forwards stream data to a given
// writable stream.
func ForwardTo(ch stream.Writable) rivers.Consumer {
	return func(upstream stream.Reader) {
		// Prevents upstream context from closing before the consumer is done
		// consuming all data.
		reader, _ := upstream.NewDownstream()

		go func() {
			defer close(ch)
			defer reader.Close(nil)

			for data := range upstream.Read() {
				select {
				case <-upstream.Done():
					return
				case ch <- data:
				}
			}
		}()
	}
}
