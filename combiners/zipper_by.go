package combiners

import "github.com/drborges/rivers/rx"

type zipBy struct {
	context rx.Context
	fn      rx.ReduceFn
}

func (c *zipBy) Combine(in ...rx.Readable) rx.Readable {
	max := func(rs ...rx.Readable) int {
		max := 0
		for _, r := range rs {
			capacity := cap(r)
			if max < capacity {
				max = capacity
			}
		}
		return max
	}

	reader, writer := rx.NewStream(max(in...))

	go func() {
		defer c.context.Recover()
		defer close(writer)

		var zipped rx.T
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
