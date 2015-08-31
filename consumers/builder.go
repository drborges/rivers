package consumers

import (
	"errors"
	"github.com/drborges/rivers/stream"
	"reflect"
)

var (
	ErrNoSuchPointer      = errors.New("Element is not a pointer")
	ErrNoSuchSlicePointer = errors.New("Element is not a pointer to a slice")
)

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (builder *Builder) Drainer() stream.Consumer {
	return &drainer{builder.context}
}

func (builder *Builder) ItemsCollector(dst interface{}) stream.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		panic(ErrNoSuchSlicePointer)
	}

	return &itemsCollector{
		context:   builder.context,
		container: slicePtr.Elem(),
	}
}

func (builder *Builder) LastItemCollector(dst interface{}) stream.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr {
		panic(ErrNoSuchPointer)
	}

	return &itemCollector{
		context: builder.context,
		item:    slicePtr.Elem(),
	}
}
