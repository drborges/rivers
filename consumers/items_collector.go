package consumers

import (
	"github.com/drborges/rivers/rx"
	"reflect"
)

type itemsCollector struct {
	context   rx.Context
	container reflect.Value
}

func (collector *itemsCollector) Consume(in rx.Readable) {
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
