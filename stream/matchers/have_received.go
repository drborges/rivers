package matchers

import (
	"fmt"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

// HaveReceived returns a matcher that verifies if the given stream.Reader has
// received the given sequence of items.
func HaveReceived(items ...interface{}) expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for _, item := range items {
			data, more := <-reader.Read()

			if !more {
				return fmt.Errorf("Expected stream to have received %v, but there was no more data", item)
			}

			if data != item {
				return fmt.Errorf("Expected %v, got %v", item, data)
			}
		}

		return nil
	}
}
