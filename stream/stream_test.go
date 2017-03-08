package stream_test

import (
	goContext "context"
	"testing"
	"time"

	"github.com/drborges/rivers/context"
	"github.com/drborges/rivers/expectations"
	. "github.com/drborges/rivers/expectations/matchers"
	"github.com/drborges/rivers/stream"
	. "github.com/drborges/rivers/stream/matchers"
)

func TestReaderWriterStreamComponents(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New()

	if err := expect(reader).To(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestReaderWriterStreamComponentsWithCustomContext(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.NewWithContext(context.New())

	if err := expect(reader).To(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestReaderDoesNotReceiveDataFromWriterWhenWriterIsClosed(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New()

	writer.Close(nil)

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestReaderDoesNotReceiveDataFromWriterWhenReaderIsClosed(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New()

	reader.Close(nil)

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestClosingWriter(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New()

	writer.Close(nil)

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}

	if err := expect(writer.Write(1)).To(Be(goContext.Canceled)); err != nil {
		t.Error(err)
	}
}

func TestWriterTimesout(t *testing.T) {
	expect := expectations.New()

	ctx := context.WithConfig(context.New(), context.Config{
		Timeout:    500 * time.Millisecond,
		BufferSize: 0,
	})

	_, writer := stream.NewWithContext(ctx)

	if err := expect(writer).To(TimeoutWithin(500 * time.Millisecond)); err != nil {
		t.Error(err)
	}
}

func TestWriterReturnsTimeoutError(t *testing.T) {
	expect := expectations.New()

	ctx := context.WithConfig(context.New(), context.Config{
		Timeout:    500 * time.Millisecond,
		BufferSize: 0,
	})

	_, writer := stream.NewWithContext(ctx)

	if err := expect(writer.Write(1)).To(Be(goContext.DeadlineExceeded)); err != nil {
		t.Error(err)
	}
}

func TestDownstream(t *testing.T) {
	expect := expectations.New()

	upstream, _ := stream.New()
	reader, writer := upstream.NewDownstream()

	if err := expect(reader).To(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}

func TestUpstreamIsClosedAfterAllDownstreamsAreClosed(t *testing.T) {
	expect := expectations.New()

	upstreamReader, upstreamWriter := stream.New()
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
