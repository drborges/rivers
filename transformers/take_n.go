package transformers

import "github.com/drborges/rivers/rx"

type takeN struct {
	context rx.Context
	n       int
}

func (take *takeN) Transform(in rx.Readable) rx.Readable {
	reader, writer := rx.NewStream(cap(in))

	go func() {
		defer take.context.Recover()
		defer close(writer)

		taken := 0

		for {
			select {
			case <-take.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				if taken >= take.n {
					return
				}

				writer <- data
				taken++
			}
		}
	}()

	return reader
}
