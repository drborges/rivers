package combiners

import (
	"github.com/drborges/rivers/stream"
	"time"
)

type zipBy struct {
	context stream.Context
	fn      stream.ReduceFn
}

func ZipBy(fn stream.ReduceFn) stream.Combiner {
	return &zipBy{
		fn: fn,
	}
}

func (combiner *zipBy) Attach(context stream.Context) {
	combiner.context = context
}

func (combiner *zipBy) Combine(in ...stream.Readable) stream.Readable {
	max := func(rs ...stream.Readable) int {
		max := 0
		for _, r := range rs {
			capacity := r.Capacity()
			if max < capacity {
				max = capacity
			}
		}
		return max
	}

	reader, writer := stream.New(max(in...))

	go func() {
		defer combiner.context.Recover()
		defer close(writer)

		var zipped stream.T
		doneIndexes := make(map[int]bool)

		for len(doneIndexes) < len(in) {
			select {
			case <-combiner.context.Failure():
				return
			case <-time.After(combiner.context.Deadline()):
				panic(stream.Timeout)
			default:
				for i, readable := range in {
					data, opened := <-readable

					if !opened {
						if _, registered := doneIndexes[i]; !registered {
							doneIndexes[i] = true
						}
						continue
					}

					if zipped == nil {
						zipped = data
					} else {
						zipped = combiner.fn(zipped, data)
					}
				}

				if zipped != nil {
					writer <- zipped
					zipped = nil
				}
			}
		}
	}()

	return reader
}
