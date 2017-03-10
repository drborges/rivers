package matchers

import (
	"fmt"

	"github.com/drborges/rivers"
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

type To func(ch stream.Readable) expectations.MatchFunc

func (fn To) To(ch stream.Readable) expectations.MatchFunc {
	return fn(ch)
}

// Filter returns a matcher that can verify if a given consumer filters the
// expected items from a sequence of items received from its upstream.
func Forward(expectedItems ...interface{}) To {
	return func(ch stream.Readable) expectations.MatchFunc {
		return func(actual interface{}) error {
			consumer, ok := actual.(rivers.Consumer)

			if !ok {
				return fmt.Errorf("Exected an actual that implements 'rivers.Consumer', got %v", actual)
			}

			reader, writer := stream.New(ctxtree.New())
			consumer(reader)

			for _, expected := range expectedItems {
				writer.Write(expected)

				if data := <-ch; data != expected {
					return fmt.Errorf("Expected consumer to have received %v, got %v", expected, data)
				}
			}

			return nil
		}
	}
}
