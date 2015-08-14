package transformers

import (
	"github.com/drborges/rivers/rx"
)

type each struct {
	context rx.Context
	handler rx.HandleFn
}

func (p *each) Transform(in rx.InStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

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
