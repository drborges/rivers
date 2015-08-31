package producers

import (
	"github.com/drborges/rivers/rx"
	"reflect"
)

type fromSlice struct {
	context rx.Context
	slice   rx.T
}

func (p *fromSlice) Produce() rx.Readable {
	sv := reflect.ValueOf(p.slice)

	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Ptr {
		panic("No such slice")
	}

	if sv.Kind() == reflect.Ptr && sv.Elem().Kind() != reflect.Slice {
		panic("No such slice")
	}

	reader, writer := rx.NewStream(sv.Len())

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
