package transformers

import "github.com/drborges/riversv2/rx"

type mapper struct {
	context rx.Context
	mapFn   rx.MapFn
}

func (t *mapper) Transform(in rx.InStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

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
