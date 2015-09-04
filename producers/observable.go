package producers

import "github.com/drborges/rivers/stream"

type Emitter struct {
	context  stream.Context
	writable stream.Writable
}

func (emitter *Emitter) Emit(data stream.T) {
	select {
	case <-emitter.context.Closed():
		panic("Context is closed")
	default:
		emitter.writable <- data
	}
}

type Observable struct {
	Context  stream.Context
	Capacity int
	Emit     func(stream.Emitter)
}

func (observable *Observable) Produce() stream.Readable {
	if observable.Capacity <= 0 {
		observable.Capacity = 10
	}
	readable, writable := stream.New(observable.Capacity)

	go func() {
		defer observable.Context.Recover()
		defer close(writable)

		if observable.Emit != nil {
			observable.Emit(&Emitter{observable.Context, writable})
		}
	}()

	return readable
}
