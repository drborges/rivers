package transformers

import (
	"github.com/drborges/rivers/rx"
)

type batch struct {
	size  int
	items []rx.T
}

func (batch *batch) Full() bool {
	return len(batch.items) == batch.size
}

func (batch *batch) Empty() bool {
	return len(batch.items) == 0
}

func (batch *batch) Commit(out rx.OutStream) {
	out <- batch.items
	batch.items = []rx.T{}
}

func (batch *batch) Add(data rx.T) {
	batch.items = append(batch.items, data)
}

type batcher struct {
	context rx.Context
	batch   rx.Batch
}

func (t *batcher) Transform(in rx.InStream) rx.InStream {
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
					if !t.batch.Empty() {
						t.batch.Commit(writer)
					}
					return
				}
				t.batch.Add(data)
				if t.batch.Full() {
					t.batch.Commit(writer)
				}
			}
		}
	}()

	return reader
}
