package matchers

import (
	"fmt"

	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/pipeline"
)

// Produce matcher that verifies if a given producer produces the given items in
// that order.
func Produce(items ...int) expectations.MatchFunc {
	return func(actual interface{}) error {
		producer, ok := actual.(pipeline.Producer)

		if !ok {
			return fmt.Errorf("Expected actual to implement 'pipeline.Producer', got %v", actual)
		}

		reader := producer(ctxtree.New())

		for _, item := range items {
			if data := <-reader.Read(); data != item {
				return fmt.Errorf("Expected producer to have produced %v, got %v", item, data)
			}
		}

		return nil
	}
}
