package aggregators_test

import (
	"testing"

	"github.com/drborges/rivers/aggregators"
	"github.com/drborges/rivers/ctxtree"
	. "github.com/drborges/rivers/ctxtree/matchers"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func TesdtFIFO(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.New()
	r1, w1 := stream.New(ctx)
	r2, w2 := stream.New(ctx)

	w1.Write(1)
	w2.Write(2, 3)

	reader := aggregators.FIFO(r1, r2)

	if err := expect(reader).To(HaveReceived(1, 2, 3)); err != nil {
		t.Error(err)
	}
}

func TestFIFOIsClosedWhenUpstreamsAreClosed(t *testing.T) {
	expect := expectations.New()

	r1, w1 := stream.New(ctxtree.New())
	r2, w2 := stream.New(ctxtree.New())

	w1.Close(nil)
	w2.Close(nil)

	reader := aggregators.FIFO(r1, r2)

	if err := expect(reader).To(BeClosed()); err != nil {
		t.Error(err)
	}
}

func TestFIFOUpstreamsAreClosedWhenDownstreamIsClosed(t *testing.T) {
	expect := expectations.New()

	r1, _ := stream.New(ctxtree.New())
	r2, _ := stream.New(ctxtree.New())
	reader := aggregators.FIFO(r1, r2)

	reader.Close(nil)

	if err := expect(r1).To(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(r2).To(BeClosed()); err != nil {
		t.Error(err)
	}
}
