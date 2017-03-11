package matchers

import (
	"fmt"

	"github.com/drborges/rivers"
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

type FromUpstream func(upstreamItems []interface{}) ToDownstreams

func (upstream FromUpstream) FromUpstream(upstreamItems ...interface{}) ToDownstreams {
	return upstream(upstreamItems)
}

type ToDownstreams struct {
	expected []interface{}
	from     []interface{}
}

func (forward ToDownstreams) ToFirstDownstream() expectations.MatchFunc {
	return func(actual interface{}) error {
		splitter, ok := actual.(rivers.Splitter)

		if !ok {
			return fmt.Errorf("Expected an actual that implements 'rivers.Splitter', got %v", actual)
		}

		reader, writer := stream.New(ctxtree.New())

		for _, item := range forward.from {
			writer.Write(item)
		}

		firstDownstream, secondDownstream := splitter(reader)

		for _, item := range forward.expected {
			data, more := <-firstDownstream.Read()

			if !more {
				return fmt.Errorf("Expected first downstream to have received %v, but there was no more data to consume", item)
			}

			if data != item {
				return fmt.Errorf("Expected first downstream to have received %v, got %v", item, data)
			}

			if data := <-secondDownstream.Read(); data == item {
				return fmt.Errorf("Expected second downstream to have not received %v", item)
			}
		}

		return nil
	}
}

func (forward ToDownstreams) ToSecondDownstream() expectations.MatchFunc {
	return func(actual interface{}) error {
		splitter, ok := actual.(rivers.Splitter)

		if !ok {
			return fmt.Errorf("Expected an actual that implements 'rivers.Splitter', got %v", actual)
		}

		reader, writer := stream.New(ctxtree.New())

		for _, item := range forward.from {
			writer.Write(item)
		}

		firstDownstream, secondDownstream := splitter(reader)

		for _, item := range forward.expected {
			data, more := <-secondDownstream.Read()

			if !more {
				return fmt.Errorf("Expected first downstream to have received %v, but there was no more data to consume", item)
			}

			if data != item {
				return fmt.Errorf("Expected first downstream to have received %v, got %v", item, data)
			}

			if data := <-firstDownstream.Read(); data == item {
				return fmt.Errorf("Expected second downstream to have not received %v", item)
			}
		}

		return nil
	}
}

func Forward(items ...interface{}) FromUpstream {
	return func(upstreamItems []interface{}) ToDownstreams {
		return ToDownstreams{
			expected: items,
			from:     upstreamItems,
		}
	}
}
