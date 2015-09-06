package combiners

import (
	"github.com/drborges/rivers/stream"
	"sync"
)

type fifo struct {
	context stream.Context
}

func FIFO() stream.Combiner {
	return &fifo{}
}

func (c *fifo) Bind(context stream.Context) {
	c.context = context
}

func (c *fifo) Combine(in ...stream.Readable) stream.Readable {
	capacity := func(in ...stream.Readable) int {
		capacity := 0
		for _, r := range in {
			capacity += cap(r)
		}
		return capacity
	}

	var wg sync.WaitGroup
	reader, writer := stream.New(capacity(in...))

	for _, r := range in {
		wg.Add(1)
		go func(r stream.Readable) {
			defer c.context.Recover()
			defer wg.Done()

			select {
			case <-c.context.Done():
				return
			default:
				for data := range r {
					writer <- data
				}
			}
		}(r)
	}

	go func() {
		defer close(writer)
		wg.Wait()
	}()

	return reader
}
