package dispatchers

import (
	"github.com/drborges/rivers/stream"
	"sync"
)

type ifDispatcher struct {
	context stream.Context
	fn      stream.PredicateFn
}

func (dispatcher *ifDispatcher) Dispatch(in stream.Readable, out ...stream.Writable) stream.Readable {
	reader, writer := stream.New(cap(in))

	wg := make(map[stream.Writable]*sync.WaitGroup)
	closeWritables := func() {
		close(writer)
		for _, writable := range out {
			go func(s stream.Writable) {
				wg[s].Wait()
				close(s)
			}(writable)
		}
	}

	go func() {
		defer dispatcher.context.Recover()
		defer closeWritables()

		for _, s := range out {
			wg[s] = &sync.WaitGroup{}
		}

		for data := range in {
			select {
			case <-dispatcher.context.Closed():
				return
			default:
				if dispatcher.fn(data) {
					for _, writable := range out {
						wg[writable].Add(1)
						// dispatch data asynchronously so that
						// slow receivers don't block the dispatch
						// process
						go func(s stream.Writable, d stream.T) {
							s <- d
							wg[s].Done()
						}(writable, data)
					}
				} else {
					writer <- data
				}
			}
		}
	}()

	return reader
}
