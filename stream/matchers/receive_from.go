package matchers

import (
	"fmt"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

// Receive allows constructing a matcher to verify whether a stream.Reader
// receives the given sequence of items From the provided stream.Writer.
func Receive(items ...interface{}) ReceiveFrom {
	return ReceiveFrom{
		items: items,
	}
}

type ReceiveFrom struct {
	items []interface{}
}

func (receive ReceiveFrom) FromWriter() expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for _, item := range receive.items {
			select {
			case data, more := <-reader.Read():
				if !more {
					return fmt.Errorf("Expected stream.Reader to have received %v from upstream, but there was no more data", item)
				}

				if data != item {
					return fmt.Errorf("Expected stream.Reader to have received %v from upstream, got %v", item, data)
				}
			default:
				return fmt.Errorf("Expected stream.Reader to have received %v from upstream, but there was no more data", item)
			}
		}

		return nil
	}
}
