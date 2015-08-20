package consumers

import "github.com/drborges/rivers/rx"

type itemsCollector struct {
	context rx.Context
	data    *[]rx.T
}

func (collector *itemsCollector) Consume(in rx.InStream) {
	for {
		select {
		case <-collector.context.Closed():
			return
		case item, more := <-in:
			if !more {
				return
			}

			*collector.data = append(*collector.data, item)
		}
	}
}
