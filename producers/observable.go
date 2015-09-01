package producers

import "github.com/drborges/rivers/stream"

type Observable struct {
	Context  stream.Context
	Capacity int
	Emit     func(w stream.Writable)
}

func (observable *Observable) Produce() stream.Readable {
	if observable.Capacity <= 0 {
		observable.Capacity = 10
	}
	reader, writer := stream.New(observable.Capacity)

	go func() {
		defer observable.Context.Recover()
		defer close(writer)

		if observable.Emit != nil {
			observable.Emit(writer)
		}
	}()

	return reader
}
