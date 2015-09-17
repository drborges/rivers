package consumers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
	"time"
)

type itemsCollector struct {
	context   stream.Context
	container reflect.Value
}

func (collector *itemsCollector) Bind(context stream.Context) {
	collector.context = context
}

func (collector *itemsCollector) Consume(in stream.Readable) {
	defer collector.context.Recover()

	for {
		select {
		case <-collector.context.Failure():
			return
		case <-time.After(collector.context.Deadline()):
			panic(stream.Timeout)
		case item, more := <-in:
			if !more {
				return
			}

			collector.container.Set(reflect.Append(collector.container, reflect.ValueOf(item)))
		}
	}
}
