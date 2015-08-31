package transformers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type flatten struct {
	context stream.Context
}

func (t *flatten) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		for {
			select {
			case <-t.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					return
				}

				dv := reflect.ValueOf(data)
				if dv.Kind() == reflect.Slice || dv.Kind() == reflect.Ptr && dv.Elem().Kind() == reflect.Slice {
					for i := 0; i < dv.Len(); i++ {
						writer <- dv.Index(i).Interface()
					}
				} else {
					writer <- data
				}
			}
		}
	}()

	return reader
}
