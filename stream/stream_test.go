package stream_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drborges/rivers/ctxtree"
	. "github.com/drborges/rivers/ctxtree/matchers"
	"github.com/drborges/rivers/expectations"
	. "github.com/drborges/rivers/expectations/matchers"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func TestReaderWriterStreamComponents(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New(ctxtree.New())

	if err := expect(reader).To(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestReaderDoesNotReceiveDataFromWriterWhenWriterIsClosed(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New(ctxtree.New())

	writer.Close(nil)

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestReaderDoesNotReceiveDataFromWriterWhenReaderIsClosed(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New(ctxtree.New())

	reader.Close(nil)
	reader.Close(nil) // accounts for multiple attempts to close a channel

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestClosingWriter(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New(ctxtree.New())

	writer.Close(nil)
	writer.Close(nil) // accounts for multiple attempts to close a channel

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}

	if err := expect(writer.Write(1)).To(Be(context.Canceled)); err != nil {
		t.Error(err)
	}
}

func TestWriterTimesout(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.WithConfig(ctxtree.New(), ctxtree.Config{
		Timeout:    50 * time.Millisecond,
		BufferSize: 0,
	})

	_, writer := stream.New(ctx)

	if err := expect(writer).To(TimeoutWithin(100 * time.Millisecond)); err != nil {
		t.Error(err)
	}
}

func TestWriterReturnsTimeoutError(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.WithConfig(ctxtree.New(), ctxtree.Config{
		Timeout:    50 * time.Millisecond,
		BufferSize: 0,
	})

	_, writer := stream.New(ctx)

	if err := expect(writer.Write(1)).To(Be(context.DeadlineExceeded)); err != nil {
		t.Error(err)
	}
}

func TestDownstream(t *testing.T) {
	expect := expectations.New()

	upstream, _ := stream.New(ctxtree.New())
	reader, writer := upstream.NewDownstream()

	if err := expect(reader).To(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestUpstreamIsClosedAfterAllDownstreamsAreClosed(t *testing.T) {
	expect := expectations.New()

	upstreamReader, upstreamWriter := stream.New(ctxtree.New())
	r1, w1 := upstreamReader.NewDownstream()
	r2, w2 := upstreamReader.NewDownstream()

	r1.Close(nil)

	if err := expect(r1).ToNot(Receive(1, 2, 4).From(w1)); err != nil {
		t.Error(err)
	}

	if err := expect(r2).To(Receive(1, 2, 4).From(w2)); err != nil {
		t.Error(err)
	}

	if err := expect(upstreamReader).To(Receive(1, 2, 4).From(upstreamWriter)); err != nil {
		t.Error(err)
	}

	r2.Close(nil)

	if err := expect(r1).ToNot(Receive(1, 2, 4).From(w1)); err != nil {
		t.Error(err)
	}

	if err := expect(r2).ToNot(Receive(1, 2, 4).From(w2)); err != nil {
		t.Error(err)
	}

	if err := expect(upstreamReader).ToNot(Receive(1, 2, 4).From(upstreamWriter)); err != nil {
		t.Error(err)
	}
}

func TestDownstreamErrorIsCollectedByRootContext(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.New()
	r1, _ := stream.New(ctx)
	r2, _ := r1.NewDownstream()

	err := errors.New("panic!")
	r2.Close(err)

	if err := expect(ctx.Err()).To(Be(err)); err != nil {
		t.Error(err)
	}
}

func TestDownstreamErrorClosesContext(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.New()
	r1, _ := stream.New(ctx)
	r2, _ := r1.NewDownstream()

	r2.Close(errors.New("panic!"))

	if err := expect(ctx).To(BeClosed()); err != nil {
		t.Error(err)
	}
}

func TestDownstreamErrorClosesAllStreams(t *testing.T) {
	expect := expectations.New()

	ctx := ctxtree.New()
	r1, w1 := stream.New(ctx)
	r2, w2 := r1.NewDownstream()
	r3, w3 := r2.NewDownstream()

	r2.Close(errors.New("panic!"))

	if err := expect(r1).ToNot(Receive(1, 2, 4).From(w1)); err != nil {
		t.Error(err)
	}

	if err := expect(r2).ToNot(Receive(1, 2, 4).From(w2)); err != nil {
		t.Error(err)
	}

	if err := expect(r3).ToNot(Receive(1, 2, 4).From(w3)); err != nil {
		t.Error(err)
	}
}
