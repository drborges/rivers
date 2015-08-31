package combiners

import "github.com/drborges/rivers/rx"

type zip struct {
	context rx.Context
}

func (c *zip) Combine(in ...rx.Readable) rx.Readable {
	capacity := func(rs ...rx.Readable) int {
		capacity := 0
		for _, r := range rs {
			capacity += cap(r)
		}
		return capacity
	}

	reader, writer := rx.NewStream(capacity(in...))

	go func() {
		defer c.context.Recover()
		defer close(writer)

		for {
			select {
			case <-c.context.Closed():
				return
			default:
				doneCount := 0
				for _, readable := range in {
					data, more := <-readable
					if !more {
						doneCount++
						continue
					}
					writer <- data
				}

				if doneCount == len(in) {
					return
				}
			}
		}
	}()

	return reader
}
