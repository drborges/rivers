package combiners

import (
	"github.com/drborges/rivers/stream"
	"time"
)

type zip struct {
	context stream.Context
}

func Zip() stream.Combiner {
	return &zip{}
}

func (combiner *zip) Attach(context stream.Context) {
	combiner.context = context
}

func (combiner *zip) Combine(in ...stream.Readable) stream.Readable {
	capacity := func(rs ...stream.Readable) int {
		capacity := 0
		for _, r := range rs {
			capacity += r.Capacity()
		}
		return capacity
	}

	reader, writer := stream.New(capacity(in...))

	go func() {
		defer combiner.context.Recover()
		defer close(writer)

		for {
			select {
			case <-combiner.context.Failure():
				return
			case <-time.After(combiner.context.Deadline()):
				panic(stream.Timeout)
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
