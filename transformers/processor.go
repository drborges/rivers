package transformers

import "github.com/drborges/rivers/stream"

type processor struct {
	context stream.Context
	onData  stream.OnDataFn
}

func (processor *processor) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

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
