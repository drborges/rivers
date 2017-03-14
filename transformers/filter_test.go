package transformers_test

import (
	"testing"

	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
	"github.com/drborges/rivers/transformers"
)

func TestFilterEvenIntegers(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	reader, writer := stream.New(ctxtree.New())
	writer.Write(1, 2, 3, 4)

	evensStream := transformers.Filter(evens)(reader)

	if err := expect(evensStream).To(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}
}

func TestDoesNotFilterStreamWhenContextIsClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	reader, writer := stream.New(ctxtree.New())
	writer.Close(nil)
	writer.Write(1, 2, 3, 4)

	evensStream := transformers.Filter(evens)(reader)

	if err := expect(evensStream).ToNot(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}
}
