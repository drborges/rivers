package matchers

import (
	"fmt"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

// From function that implements the "From" part of the "Receive/From" mathcer,
// returning the actual matcher.
type From func(stream.Writer) expectations.MatchFunc

// From provides the "From" part of the Receive(sequence).From(stream.Writer)
// matcher DSL.
func (fn From) From(writer stream.Writer) expectations.MatchFunc {
	return fn(writer)
}

// Receive allows constructing a matcher to verify whether a stream.Reader
// receives the given sequence of items From the provided stream.Writer.
func Receive(items ...interface{}) From {
	return func(writer stream.Writer) expectations.MatchFunc {
		return func(actual interface{}) error {
			reader, ok := actual.(stream.Reader)

			if !ok {
				return fmt.Errorf("Exected an actual that implements 'stream.Reader', got %v", actual)
			}

			for _, item := range items {
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
}
