package consumers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type itemsCollector struct {
	context   stream.Context
	container reflect.Value
}

func (collector *itemsCollector) Consume(in stream.Readable) {
	for {
		select {
		case <-collector.context.Closed():
			return
		case item, more := <-in:
			if !more {
				return
			}

			collector.container.Set(reflect.Append(collector.container, reflect.ValueOf(item)))
		}
	}
}
