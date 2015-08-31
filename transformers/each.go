package transformers

import (
	"github.com/drborges/rivers/stream"
)

type each struct {
	context stream.Context
	handler stream.EachFn
}

func (p *each) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer p.context.Recover()
		defer close(writer)

		for {
			select {
			case <-p.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}
				p.handler(data)
				writer <- data
			}
		}
	}()

	return reader
}
