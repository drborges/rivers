package splitters_test

import (
	"testing"

	"github.com/drborges/rivers/ctxtree"
	. "github.com/drborges/rivers/ctxtree/matchers"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/splitters"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func TestSplitByPredicate(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	reader, writer := stream.New(ctxtree.New())
	writer.Write(1, 2, 3, 4)

	evensStream, oddsStream := splitters.ByPredicate(evens)(reader)

	if err := expect(evensStream).To(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}

	if err := expect(oddsStream).To(HaveReceived(1, 3)); err != nil {
		t.Error(err)
	}
}

func TestSplitByPredicateDoesNotConsumeDataWhenUpstreamIsClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	reader, writer := stream.New(ctxtree.New())
	writer.Close(nil)
	writer.Write(1, 2, 3, 4)

	evensStream, oddsStream := splitters.ByPredicate(evens)(reader)

	if err := expect(evensStream).ToNot(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}

	if err := expect(oddsStream).ToNot(HaveReceived(1, 3)); err != nil {
		t.Error(err)
	}
}

func TestSplitByPredicateUpstreamIsClosedWhenAllDownstreamsAreClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	upstream, _ := stream.New(ctxtree.New())
	downstream1, downstream2 := splitters.ByPredicate(evens)(upstream)

	downstream1.Close(nil)

	if err := expect(upstream).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	downstream2.Close(nil)

	if err := expect(upstream).To(BeClosed()); err != nil {
		t.Error(err)
	}
}
