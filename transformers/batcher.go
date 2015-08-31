package transformers

import (
	"github.com/drborges/rivers/stream"
)

type batch struct {
	size  int
	items []stream.T
}

func (batch *batch) Full() bool {
	return len(batch.items) == batch.size
}

func (batch *batch) Empty() bool {
	return len(batch.items) == 0
}

func (batch *batch) Commit(out stream.Writable) {
	out <- batch.items
	batch.items = []stream.T{}
}

func (batch *batch) Add(data stream.T) {
	batch.items = append(batch.items, data)
}

type batcher struct {
	context stream.Context
	batch   stream.Batch
}

func (t *batcher) Transform(in stream.Readable) stream.Readable {
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
