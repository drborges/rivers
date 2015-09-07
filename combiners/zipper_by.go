package combiners

import (
	"github.com/drborges/rivers/stream"
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

func (c *zipBy) Bind(context stream.Context) {
	c.context = context
}

func (c *zipBy) Combine(in ...stream.Readable) stream.Readable {
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
		defer c.context.Recover()
		defer close(writer)

		var zipped stream.T
		doneIndexes := make(map[int]bool)

		for len(doneIndexes) < len(in) {
			select {
			case <-c.context.Done():
				return
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
						zipped = c.fn(zipped, data)
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
