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
	return &Sink{}
}

func ItemsCollector(dst interface{}) stream.Consumer {
	ptr := reflect.ValueOf(dst)

	if ptr.Kind() != reflect.Ptr || ptr.Elem().Kind() != reflect.Slice {
		panic(ErrNoSuchSlicePointer)
	}

	container := ptr.Elem()
	return &Sink{
		OnNext: func(data stream.T) {
			container.Set(reflect.Append(container, reflect.ValueOf(data)))
		},
	}
}

func LastItemCollector(dst interface{}) stream.Consumer {
	ptr := reflect.ValueOf(dst)

	if ptr.Kind() != reflect.Ptr {
		panic(ErrNoSuchPointer)
	}

	val := ptr.Elem()
	return &Sink{
		OnNext: func(data stream.T) {
			val.Set(reflect.ValueOf(data))
		},
	}
}

func CollectBy(fn stream.EachFn) stream.Consumer {
	return &Sink{
		OnNext: fn,
	}
}

func GroupBy(fn stream.MapFn, result stream.Groups) stream.Consumer {
	return &Sink{
		OnNext: func(data stream.T) {
			groupKey := fn(data)
			if _, exists := result[groupKey]; !exists {
				result[groupKey] = []stream.T{}
			}

			result[groupKey] = append(result[groupKey], data)
		},
	}
}
