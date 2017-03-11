package transformers

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

// Filter implements a rivers.Transformer which filters stream data that match
// the given predicate.
func Filter(fn rivers.Predicate) rivers.Transformer {
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
