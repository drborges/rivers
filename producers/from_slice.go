package producers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type fromSlice struct {
	context stream.Context
	slice   stream.T
}

func (p *fromSlice) Produce() stream.Readable {
	sv := reflect.ValueOf(p.slice)

	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Ptr {
		panic("No such slice")
	}

	if sv.Kind() == reflect.Ptr && sv.Elem().Kind() != reflect.Slice {
		panic("No such slice")
	}

	reader, writer := stream.New(sv.Len())

	go func() {
		defer close(writer)
		defer p.context.Recover()

		for i := 0; i < sv.Len(); i++ {
			select {
			case <-p.context.Closed():
				break
			default:
				writer <- sv.Index(i).Interface()
			}
		}
	}()

	return reader
}
