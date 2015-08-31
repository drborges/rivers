package transformers

import "github.com/drborges/rivers/stream"

type filter struct {
	context   stream.Context
	predicate stream.PredicateFn
}

func (t *filter) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		for {
			select {
			case <-t.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				if t.predicate(data) {
					writer <- data
				}
			}
		}
	}()

	return reader
}
