package aggregators

import (
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/stream"
)

func FIFO(upstream1, upstream2 stream.Reader) stream.Reader {
	reader, writer := stream.New(ctxtree.New())

	upstream1Done := make(chan struct{}, 0)
	upstream2Done := make(chan struct{}, 0)

	go func() {
		defer close(upstream1Done)
		defer upstream1.Close(nil)

		for data := range upstream1.Read() {
			if err := writer.Write(data); err != nil {
				return
			}
		}
	}()

	go func() {
		defer close(upstream2Done)
		defer upstream2.Close(nil)

		for data := range upstream2.Read() {
			if err := writer.Write(data); err != nil {
				return
			}
		}
	}()

	go func() {
		defer writer.Close(nil)
		<-upstream1Done
		<-upstream2Done
	}()

	return reader
}
