package transformers

import "github.com/drborges/rivers/stream"

type Observer struct {
	Context     stream.Context
	OnCompleted func(w stream.Writable)
	OnNext      func(data stream.T, w stream.Writable) error
}

func (transformer *Observer) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer transformer.Context.Recover()
		defer close(writer)

		for {
			select {
			case <-transformer.Context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					if transformer.OnCompleted != nil {
						transformer.OnCompleted(writer)
					}
					return
				}

				if transformer.OnNext == nil {
					continue
				}

				if err := transformer.OnNext(data, writer); err != nil {
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
