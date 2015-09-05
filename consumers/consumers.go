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

func Drainer() stream.Consumer {
	return &drainer{}
}

func ItemsCollector(dst interface{}) stream.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr || slicePtr.Elem().Kind() != reflect.Slice {
		panic(ErrNoSuchSlicePointer)
	}

	return &itemsCollector{
		container: slicePtr.Elem(),
	}
}

func LastItemCollector(dst interface{}) stream.Consumer {
	slicePtr := reflect.ValueOf(dst)

	if slicePtr.Kind() != reflect.Ptr {
		panic(ErrNoSuchPointer)
	}

	return &itemCollector{
		item: slicePtr.Elem(),
	}
}
