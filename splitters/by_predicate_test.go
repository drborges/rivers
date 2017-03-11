package splitters_test

import (
	"testing"

	"github.com/drborges/rivers/ctxtree"
	. "github.com/drborges/rivers/ctxtree/matchers"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/splitters"
	. "github.com/drborges/rivers/splitters/matchers"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func TestSplitByPredicate(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	splitter := splitters.ByPredicate(evens)

	if err := expect(splitter).To(Forward(0, 2).FromUpstream(0, 1, 2, 3).ToFirstDownstream()); err != nil {
		t.Error(err)
	}

	if err := expect(splitter).To(Forward(1, 3).FromUpstream(0, 1, 2, 3).ToSecondDownstream()); err != nil {
		t.Error(err)
	}
}

func TestSplitByPredicateDoesNotConsumeDataWhenUpstreamIsClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	splitter := splitters.SplitterWithClosedUpstream(splitters.ByPredicate(evens))

	if err := expect(splitter).ToNot(Forward(0, 2).FromUpstream(0, 1, 2, 3).ToFirstDownstream()); err != nil {
		t.Error(err)
	}

	if err := expect(splitter).ToNot(Forward(1, 3).FromUpstream(0, 1, 2, 3).ToSecondDownstream()); err != nil {
		t.Error(err)
	}
}

func TestSplitByPredicateUpstreamIsClosedWhenAllDownstreamsAreClosed(t *testing.T) {
	expect := expectations.New()
	evens := func(data stream.T) bool { return data.(int)%2 == 0 }

	ctx := ctxtree.New()
	reader, writer := stream.New(ctx)
	downstream1, downstream2 := splitters.ByPredicate(evens)(reader)

	downstream1.Close(nil)

	writer.Write(1)
	writer.Write(2)
	writer.Write(3)
	writer.Write(4)

	if err := expect(ctx).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(downstream1).ToNot(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}

	if err := expect(downstream2).To(HaveReceived(1, 3)); err != nil {
		t.Error(err)
	}

	downstream2.Close(nil)

	writer.Write(1)
	writer.Write(2)
	writer.Write(3)
	writer.Write(4)

	if err := expect(ctx).To(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(downstream1).ToNot(HaveReceived(2, 4)); err != nil {
		t.Error(err)
	}

	if err := expect(downstream2).ToNot(HaveReceived(1, 3)); err != nil {
		t.Error(err)
	}
}
