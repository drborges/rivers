package combiners

import (
	"github.com/drborges/rivers/stream"
	"sync"
	"time"
)

type fifo struct {
	context stream.Context
}

func FIFO() stream.Combiner {
	return &fifo{}
}

func (combiner *fifo) Attach(context stream.Context) {
	combiner.context = context
}

func (combiner *fifo) Combine(in ...stream.Readable) stream.Readable {
	capacity := func(in ...stream.Readable) int {
		capacity := 0
		for _, r := range in {
			capacity += r.Capacity()
		}
		return capacity
	}

	var wg sync.WaitGroup
	reader, writer := stream.New(capacity(in...))

	for _, r := range in {
		wg.Add(1)
		go func(r stream.Readable) {
			defer combiner.context.Recover()
			defer wg.Done()

			select {
			case <-combiner.context.Failure():
				return
			case <-time.After(combiner.context.Deadline()):
				panic(stream.Timeout)
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
