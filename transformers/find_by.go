package transformers

import "github.com/drborges/rivers/rx"

type findBy struct {
	context   rx.Context
	predicate rx.PredicateFn
}

func (t *findBy) Transform(in rx.Readable) rx.Readable {
	reader, writer := rx.NewStream(cap(in))

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
					return
				}
			}
		}
	}()

	return reader
}
