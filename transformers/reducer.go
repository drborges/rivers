package transformers

import "github.com/drborges/rivers/stream"

type reducer struct {
	context    stream.Context
	reduceFn   stream.ReduceFn
	initialAcc stream.T
}

func (t *reducer) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

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
