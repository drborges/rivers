package matchers

import (
	"fmt"

	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

type From func([]interface{}) expectations.MatchFunc

func (fn From) From(streamItems ...interface{}) expectations.MatchFunc {
	return fn(streamItems)
}

// Filter returns a matcher that can verify if a given transformer filters the
// expected items from a sequence of items received from its upstream.
func Filter(expectedItems ...interface{}) From {
	return func(streamItems []interface{}) expectations.MatchFunc {
		return func(actual interface{}) error {
			transformer, ok := actual.(pipeline.Transformer)

			if !ok {
				return fmt.Errorf("Exected an actual that implements 'pipeline.Transformer', got %v", actual)
			}

			reader, writer := stream.New(ctxtree.New())
			filtered := transformer(reader)

			go func() {
				for _, i := range streamItems {
					writer.Write(i)
				}
			}()

			for _, expected := range expectedItems {
				if data := <-filtered.Read(); data != expected {
					return fmt.Errorf("Expected transformer to have received %v, got %v", expected, data)
				}
			}

			return nil
		}
	}
}
