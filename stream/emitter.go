package stream

import "time"

type emitter struct {
	context  Context
	writable Writable
}

func NewEmitter(context Context, w Writable) Emitter {
	return &emitter{context, w}
}

func (emitter *emitter) Emit(data T) {
	select {
	case <-emitter.context.Done():
		panic(Done)
	case <-emitter.context.Failure():
		panic(Done)
	case <-time.After(emitter.context.Deadline()):
		panic(Timeout)
	default:
		emitter.writable <- data
	}
}
