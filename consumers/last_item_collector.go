package consumers

import (
	"github.com/drborges/rivers/rx"
	"reflect"
)

type itemCollector struct {
	context rx.Context
	item    reflect.Value
}

func (collector *itemCollector) Consume(in rx.InStream) {
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
