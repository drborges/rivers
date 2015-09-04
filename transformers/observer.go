package transformers

import "github.com/drborges/rivers/stream"

type Observer struct {
	context     stream.Context
	OnCompleted func(w stream.Writable)
	OnNext      func(data stream.T, w stream.Writable) error
}

func (observer *Observer) Bind(context stream.Context) {
	observer.context = context
}

func (observer *Observer) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer observer.context.Recover()
		defer close(writer)

		for {
			select {
			case <-observer.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					if observer.OnCompleted != nil {
						observer.OnCompleted(writer)
					}
					return
				}

				if observer.OnNext == nil {
					continue
				}

				if err := observer.OnNext(data, writer); err != nil {
					if err == stream.Done {
						return
					}
					panic(err)
				}
			}
		}
	}()

	return reader
}
