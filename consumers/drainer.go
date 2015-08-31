package consumers

import "github.com/drborges/rivers/stream"

type drainer struct {
	context stream.Context
}

func (drainer *drainer) Consume(in stream.Readable) {
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
