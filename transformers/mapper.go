package transformers

import "github.com/drborges/rivers/stream"

type mapper struct {
	context stream.Context
	mapFn   stream.MapFn
}

func (t *mapper) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		for {
			select {
			case <-t.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				writer <- t.mapFn(data)
			}
		}
	}()

	return reader
}
