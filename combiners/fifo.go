package combiners

import (
	"github.com/drborges/riversv2/rx"
	"sync"
)

type fifo struct {
	context rx.Context
}

func (c *fifo) Combine(in ...rx.InStream) rx.InStream {
	capacity := func(in ...rx.InStream) int {
		capacity := 0
		for _, r := range in {
			capacity += cap(r)
		}
		return capacity
	}

	var wg sync.WaitGroup
	reader, writer := rx.NewStream(capacity(in...))

	for _, r := range in {
		wg.Add(1)
		go func(r rx.InStream) {
			defer c.context.Recover()
			defer wg.Done()

			select {
			case <-c.context.Closed():
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
