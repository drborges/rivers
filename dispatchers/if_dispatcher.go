package dispatchers

import (
	"github.com/drborges/rivers/rx"
	"sync"
)

type ifDispatcher struct {
	context rx.Context
	fn      rx.PredicateFn
}

func (dispatcher *ifDispatcher) Dispatch(in rx.Readable, out ...rx.Writable) rx.Readable {
	reader, writer := rx.NewStream(cap(in))

	wg := make(map[rx.Writable]*sync.WaitGroup)
	closeToStreams := func() {
		close(writer)
		for _, writable := range out {
			go func(s rx.Writable) {
				wg[s].Wait()
				close(s)
			}(writable)
		}
	}

	go func() {
		defer dispatcher.context.Recover()
		defer closeToStreams()

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
						go func(s rx.Writable, d rx.T) {
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
