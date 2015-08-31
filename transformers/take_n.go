package transformers

import "github.com/drborges/rivers/stream"

type takeN struct {
	context stream.Context
	n       int
}

func (take *takeN) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer take.context.Recover()
		defer close(writer)

		taken := 0

		for {
			select {
			case <-take.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				if taken >= take.n {
					return
				}

				writer <- data
				taken++
			}
		}
	}()

	return reader
}
