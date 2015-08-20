package consumers

import (
	"github.com/drborges/rivers/rx"
	"reflect"
)

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (builder *Builder) Drainer() rx.Consumer {
	return &drainer{builder.context}
}

func (builder *Builder) ItemsCollector(dst interface{}) rx.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		panic(rx.ErrNoSuchSlicePointer)
	}

	return &itemsCollector{
		context:   builder.context,
		container: slicePtr.Elem(),
	}
}

func (builder *Builder) LastItemCollector(dst interface{}) rx.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr {
		panic(rx.ErrNoSuchPointer)
	}

	return &itemCollector{
		context: builder.context,
		item:    slicePtr.Elem(),
	}
}
