package matchers

import (
	"errors"
	"fmt"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

type receiveFromMatcher func(stream.Writer) expectations.MatchFunc

func (fn receiveFromMatcher) From(writer stream.Writer) expectations.MatchFunc {
	return fn(writer)
}

func Receive(items ...int) receiveFromMatcher {
	return func(writer stream.Writer) expectations.MatchFunc {
		return func(actual interface{}) error {
			reader, ok := actual.(stream.Reader)

			if !ok {
				return errors.New(fmt.Sprintf("Exected an actual that implements 'stream.Reader', got %v", actual))
			}

			for _, num := range items {
				writer.Write(num)

				select {
				case data := <-reader.Read():
					if data != num {
						return errors.New(fmt.Sprintf("Expected %v, got %v", num, data))
					}
				default:
					return errors.New(fmt.Sprintf("Expected stream to have received %v, but it was closed.", num))
				}
			}

			return nil
		}
	}
}
