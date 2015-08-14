package dispatchers

import (
	"github.com/drborges/rivers/rx"
	"sync"
)

type ifDispatcher struct {
	context rx.Context
	fn      rx.PredicateFn
}

func (dispatcher *ifDispatcher) Dispatch(in rx.InStream, out ...rx.OutStream) rx.InStream {
	reader, writer := rx.NewStream(cap(in))

	wg := make(map[rx.OutStream]*sync.WaitGroup)
	closeToStreams := func() {
		close(writer)
		for _, stream := range out {
			go func(s rx.OutStream) {
				wg[s].Wait()
				close(s)
			}(stream)
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
					for _, stream := range out {
						wg[stream].Add(1)
						// dispatch data asynchronously so that
						// slow receivers don't block the dispatch
						// process
						go func(s rx.OutStream, d rx.T) {
							s <- d
							wg[s].Done()
						}(stream, data)
					}
				} else {
					writer <- data
				}
			}
		}
	}()

	return reader
}
