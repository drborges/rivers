package aggregators

import "github.com/drborges/rivers/stream"

func FIFO(upstream1, upstream2 stream.Reader) stream.Reader {
	reader, writer := upstream1.NewDownstream()
	upstream1Done := make(chan struct{}, 0)
	upstream2Done := make(chan struct{}, 0)

	go func() {
		defer upstream1.Close(nil)
		defer close(upstream1Done)

		for data := range upstream1.Read() {
			if err := writer.Write(data); err != nil {
				return
			}
		}
	}()

	go func() {
		defer upstream2.Close(nil)
		defer close(upstream2Done)

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
