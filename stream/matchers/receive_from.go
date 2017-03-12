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

func (receive ReceiveFrom) FromWriter(writer stream.Writer) expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for _, item := range receive.items {
			writer.Write(item)

			select {
			case data := <-reader.Read():
				if data != item {
					return fmt.Errorf("Expected %v, got %v", item, data)
				}
			default:
				return fmt.Errorf("Expected stream to have received %v, but it was closed", item)
			}
		}

		return nil
	}
}

func (receive ReceiveFrom) FromWriters(writers ...stream.Writer) expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for i, item := range receive.items {
			writers[i].Write(item)

			data, more := <-reader.Read()

			if !more {
				return fmt.Errorf("Expected stream.Reader to have received %v from writer at position %v, but there was no more data", item, i)
			}

			if data != item {
				return fmt.Errorf("Expected stream.Reader to have received %v from writer at position %v, got %v", item, i, data)
			}
		}

		return nil
	}
}

func (receive ReceiveFrom) FromUpstream() expectations.MatchFunc {
	return func(actual interface{}) error {
		reader, ok := actual.(stream.Reader)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
		}

		for _, item := range receive.items {
			data, more := <-reader.Read()

			if !more {
				return fmt.Errorf("Expected stream.Reader to have received %v from upstream, but there was no more data", item)
			}

			if data != item {
				return fmt.Errorf("Expected stream.Reader to have received %v from upstream, got %v", item, data)
			}
		}

		return nil
	}
}
