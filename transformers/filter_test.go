package transformers_test

import (
	"testing"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"

	. "github.com/drborges/rivers/transformers/matchers"
)

func TestFilterEvenIntegers(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	if err := expect(transformers.Filter(evens)).To(Filter(2, 4).From(1, 2, 3, 4)); err != nil {
		t.Error(err)
	}
}

func TestDoesNotFilterStreamWhenContextIsClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }
	transformer := transformers.TransformerWithClosedUpstream(transformers.Filter(evens))

	if err := expect(transformer).ToNot(Filter(2, 4).From(1, 2, 3, 4)); err != nil {
		t.Error(err)
	}
}
