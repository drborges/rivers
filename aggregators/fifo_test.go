package aggregators_test

import (
	"testing"

	"github.com/drborges/rivers/aggregators"
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func FIFO(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.New()
	r1, w1 := stream.New(ctx)
	r2, w2 := stream.New(ctx)

	reader := aggregators.FIFO(r1, r2)

	if err := expect(reader).To(Receive(1, 2, 3).FromWriters(w1, w2, w1)); err != nil {
		t.Error(err)
	}
}

func TestFIFODoesNotReceiveFromClosedUpstream(t *testing.T) {
	expect := expectations.New()

	r1, w1 := stream.New(ctxtree.New())
	r2, w2 := stream.New(ctxtree.New())

	reader := aggregators.FIFO(r1, r2)

	r1.Close(nil)

	w1.Write(1)
	w2.Write(2)

	if err := expect(reader).To(HaveReceived(2)); err != nil {
		t.Error(err)
	}

	r2.Close(nil)

	w1.Write(1)
	w2.Write(2)

	if err := expect(reader).ToNot(HaveReceived(1)); err != nil {
		t.Error(err)
	}

	if err := expect(reader).ToNot(HaveReceived(2)); err != nil {
		t.Error(err)
	}
}
