package consumers

import "github.com/drborges/riversv2/rx"

type drainer struct {
	context rx.Context
}

func (drainer *drainer) Consume(in rx.InStream) {
	for {
		select {
		case <-drainer.context.Closed():
			return
		default:
			if _, more := <-in; !more {
				return
			}
		}
	}
}
