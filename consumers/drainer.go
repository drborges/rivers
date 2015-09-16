package consumers

import "github.com/drborges/rivers/stream"

type drainer struct {
	context stream.Context
}

func (drainer *drainer) Bind(context stream.Context) {
	drainer.context = context
}

func (drainer *drainer) Consume(in stream.Readable) {
	for {
		select {
		case <-drainer.context.Failure():
			return
		default:
			if _, more := <-in; !more {
				return
			}
		}
	}
}
