package transformers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// Predicate is a function that given an input it returns true if the input matches
// the predicate, false otherwise.
type Predicate func(stream.T) bool

// Filter implements a rivers.Transformer which filters stream data that match
// the given predicate.
func Filter(fn Predicate) rivers.Transformer {
	return func(upstream stream.Reader) stream.Reader {
		reader, writer := upstream.NewDownstream()

		go func() {
			defer writer.Close(nil)

			for data := range upstream.Read() {
				if fn(data) {
					if err := writer.Write(data); err != nil {
						return
					}
				}
			}
		}()

		return reader
	}
}
