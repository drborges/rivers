package combiners

import "github.com/drborges/rivers/stream"

type zipBy struct {
	context stream.Context
	fn      stream.ReduceFn
}

func (c *zipBy) Combine(in ...stream.Readable) stream.Readable {
	max := func(rs ...stream.Readable) int {
		max := 0
		for _, r := range rs {
			capacity := cap(r)
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
		var doneCount int

		for doneCount < len(in) {
			select {
			case <-c.context.Closed():
				return
			default:
				for _, readable := range in {
					data, opened := <-readable

					if !opened {
						doneCount++
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
