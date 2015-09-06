package consumers

import "github.com/drborges/rivers/stream"

type collectBy struct {
	context stream.Context
	fn      stream.EachFn
}

func (collector *collectBy) Bind(context stream.Context) {
	collector.context = context
}

func (collector *collectBy) Consume(in stream.Readable) {
	for {
		select {
		case <-collector.context.Done():
			return
		default:
			data, more := <-in
			if !more {
				return
			}
			collector.fn(data)
		}
	}
}
