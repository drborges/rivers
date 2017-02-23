package stream_test

import (
	"testing"

	"github.com/drborges/rivers/stream"
)

var expected = []int{1, 2, 4}

func TestReadableWritableStreams(t *testing.T) {
	reader, writer := stream.New()

	for _, num := range expected {
		writer.Write(num)

		if data := <-reader.Read(); data != num {
			t.Errorf("Expected %v, got %v", num, data)
		}
	}
}
