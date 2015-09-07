package combiners

import "github.com/drborges/rivers/stream"

type zip struct {
	context stream.Context
}

func Zip() stream.Combiner {
	return &zip{}
}

func (c *zip) Bind(context stream.Context) {
	c.context = context
}

func (c *zip) Combine(in ...stream.Readable) stream.Readable {
	capacity := func(rs ...stream.Readable) int {
		capacity := 0
		for _, r := range rs {
			capacity += r.Capacity()
		}
		return capacity
	}

	reader, writer := stream.New(capacity(in...))

	go func() {
		defer c.context.Recover()
		defer close(writer)

		for {
			select {
			case <-c.context.Done():
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
