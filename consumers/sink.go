package consumers

import (
	"github.com/drborges/rivers/stream"
	"time"
)

type Sink struct {
	context stream.Context
	OnNext  stream.EachFn
}

func (sink *Sink) Attach(context stream.Context) {
	sink.context = context
}

func (sink *Sink) Consume(in stream.Readable) {
	defer sink.context.Recover()

	for {
		select {
		case <-sink.context.Failure():
			return
		case <-time.After(sink.context.Deadline()):
			panic(stream.Timeout)
		case data, more := <-in:
			if !more {
				return
			}
			if sink.OnNext != nil {
				sink.OnNext(data)
			}
		}
	}
}
