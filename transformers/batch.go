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

func (batch *batch) Commit(emitter stream.Emitter) {
	emitter.Emit(batch.items)
	batch.items = []stream.T{}
}

func (batch *batch) Add(data stream.T) {
	batch.items = append(batch.items, data)
}
