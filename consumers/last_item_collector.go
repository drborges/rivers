package consumers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type itemCollector struct {
	context stream.Context
	item    reflect.Value
}

func (collector *itemCollector) Bind(context stream.Context) {
	collector.context = context
}

func (collector *itemCollector) Consume(in stream.Readable) {
	for {
		select {
		case <-collector.context.Done():
			return
		case item, more := <-in:
			if !more {
				return
			}

			collector.item.Set(reflect.ValueOf(item))
		}
	}
}
