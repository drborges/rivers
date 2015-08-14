package transformers

import "github.com/drborges/rivers/rx"

type processor struct {
	context rx.Context
	onData  rx.OnDataFn
}

func (processor *processor) Transform(in rx.InStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

	go func() {
		defer processor.context.Recover()
		defer close(writer)

		for {
			select {
			case <-processor.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				processor.onData(data, writer)
			}
		}
	}()

	return reader
}
