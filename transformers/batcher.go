package transformers

import (
	"github.com/drborges/riversv2/rx"
)

type batcher struct {
	context rx.Context
	size    int
}

func (t *batcher) Transform(in rx.InStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		batch := []rx.T{}
		for {
			select {
			case <-t.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					if len(batch) > 0 {
						writer <- batch
					}
					return
				}
				batch = append(batch, data)
				if len(batch) == t.size {
					writer <- batch
					batch = []rx.T{}
				}
			}
		}
	}()

	return reader
}