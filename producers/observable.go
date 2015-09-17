package producers

import "github.com/drborges/rivers/stream"

type Observable struct {
	context  stream.Context
	Capacity int
	Emit     func(stream.Emitter)
}

func (observable *Observable) Attach(context stream.Context) {
	observable.context = context
}

func (observable *Observable) Produce() stream.Readable {
	if observable.Capacity <= 0 {
		observable.Capacity = 10
	}
	readable, writable := stream.New(observable.Capacity)

	go func() {
		defer observable.context.Recover()
		defer close(writable)

		if observable.Emit != nil {
			observable.Emit(stream.NewEmitter(observable.context, writable))
		}
	}()

	return readable
}
