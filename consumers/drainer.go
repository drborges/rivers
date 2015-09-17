package consumers

import (
	"github.com/drborges/rivers/stream"
	"time"
)

type drainer struct {
	context stream.Context
}

func (drainer *drainer) Bind(context stream.Context) {
	drainer.context = context
}

func (drainer *drainer) Consume(in stream.Readable) {
	defer drainer.context.Recover()

	for {
		select {
		case <-drainer.context.Failure():
			return
		case <-time.After(drainer.context.Deadline()):
			panic(stream.Timeout)
		default:
			if _, more := <-in; !more {
				return
			}
		}
	}
}
