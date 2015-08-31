package consumers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type itemCollector struct {
	context stream.Context
	item    reflect.Value
}

func (collector *itemCollector) Consume(in stream.Readable) {
	for {
		select {
		case <-collector.context.Closed():
			return
		case item, more := <-in:
			if !more {
				return
			}

			collector.item.Set(reflect.ValueOf(item))
		}
	}
}
