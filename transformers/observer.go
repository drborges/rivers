package transformers

import "github.com/drborges/rivers/stream"

type Observer struct {
	context     stream.Context
	OnCompleted func(emitter stream.Emitter)
	OnNext      func(data stream.T, emitter stream.Emitter) error
}

func (observer *Observer) Bind(context stream.Context) {
	observer.context = context
}

func (observer *Observer) Transform(in stream.Readable) stream.Readable {
	readable, writable := stream.New(cap(in))
	emitter := stream.NewEmitter(observer.context, writable)

	go func() {
		defer observer.context.Recover()
		defer close(writable)

		for {
			select {
			case <-observer.context.Done():
				return
			default:
				data, more := <-in
				if !more {
					if observer.OnCompleted != nil {
						observer.OnCompleted(emitter)
					}
					return
				}

				if observer.OnNext == nil {
					continue
				}

				if err := observer.OnNext(data, emitter); err != nil {
					if err == stream.Done {
						return
					}
					panic(err)
				}
			}
		}
	}()

	return readable
}
