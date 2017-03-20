package rivers_test

import (
	"testing"

	"github.com/drborges/rivers"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
	"github.com/drborges/rivers/transformers"
)

func TestRiversAPI(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New(ctxtree.New())
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	rivers.
		From(producers.Range(1, 5)).
		Apply(transformers.Filter(evens)).
		Then(consumers.ForwardTo(writer.Writable())).
		Run()

	if err := expect(reader).To(Receive(2, 4)); err != nil {
		t.Error(err)
	}
}
