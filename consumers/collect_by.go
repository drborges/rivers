package consumers

import (
	"github.com/drborges/rivers/stream"
	"time"
)

type collectBy struct {
	context stream.Context
	fn      stream.EachFn
}

func (collector *collectBy) Bind(context stream.Context) {
	collector.context = context
}

func (collector *collectBy) Consume(in stream.Readable) {
	defer collector.context.Recover()

	for {
		select {
		case <-collector.context.Failure():
			return
		case <-time.After(collector.context.Deadline()):
			panic(stream.Timeout)
		default:
			data, more := <-in
			if !more {
				return
			}
			collector.fn(data)
		}
	}
}
