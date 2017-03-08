package stream_test

import (
	"testing"

	"github.com/drborges/rivers/context"
	"github.com/drborges/rivers/expectations"
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
}

func TestWriterTimesout(t *testing.T) {
	expect := expectations.New()

	reader, writer := stream.New()

	writer.Close(nil)

	if err := expect(reader).ToNot(Receive(1, 2, 4).From(writer)); err != nil {
		t.Error(err)
	}
}
