package stream

type emitter struct {
	context  Context
	writable Writable
}

func NewEmitter(context Context, w Writable) Emitter {
	return &emitter{context, w}
}

func (emitter *emitter) Emit(data T) {
	select {
	case <-emitter.context.Closed():
		panic("Context is closed")
	default:
		emitter.writable <- data
	}
}
