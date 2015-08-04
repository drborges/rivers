package transformers

import "github.com/drborges/riversv2/rx"

type reducer struct {
	context    rx.Context
	reduceFn   rx.ReduceFn
	initialAcc rx.T
}

func (t *reducer) Transform(in rx.InStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

	acc := t.initialAcc

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
					goto done
				}
				acc = t.reduceFn(acc, data)
			}
		}
	done:
		writer <- acc
	}()

	return reader
}
