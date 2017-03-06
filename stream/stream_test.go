package stream_test

import (
	"testing"

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
