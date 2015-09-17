package consumers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
	"time"
)

type itemCollector struct {
	context stream.Context
	item    reflect.Value
}

func (collector *itemCollector) Bind(context stream.Context) {
	collector.context = context
}

func (collector *itemCollector) Consume(in stream.Readable) {
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

			collector.item.Set(reflect.ValueOf(item))
		}
	}
}
