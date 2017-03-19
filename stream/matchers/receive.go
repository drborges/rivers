package matchers

import (
	"fmt"
	"time"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

// Receive returns a matcher that verifies if the given stream.Reader has
// received the given sequence of items.
func Receive(items ...interface{}) expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for _, item := range items {
			select {
			case data := <-reader.Read():
				if data != item {
					return fmt.Errorf("Expected stream.Reader to have received %v from upstream, got %v", item, data)
				}
			case <-time.After(10 * time.Millisecond):
				return fmt.Errorf("Expected stream.Reader to have received %v from upstream, but there was no more data", item)
			}
		}

		return nil
	}
}
