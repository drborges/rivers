package transformers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

func Filter(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			if fn(data) {
				emitter.Emit(data)
			}
			return nil
		},
	}
}

func FindBy(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			if fn(data) {
				emitter.Emit(data)
				return stream.Done
			}
			return nil
		},
	}
}

func TakeFirst(n int) stream.Transformer {
	taken := 0
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			if taken >= n {
				return stream.Done
			}

			emitter.Emit(data)
			taken++
			return nil
		},
	}
}

func DropFirst(n int) stream.Transformer {
	dropped := 0
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			if dropped < n {
				dropped++
				return nil
			}

			emitter.Emit(data)
			return nil
		},
	}
}

func Map(fn stream.MapFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			emitter.Emit(fn(data))
			return nil
		},
	}
}

func OnData(fn stream.OnDataFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			fn(data, emitter)
			return nil
		},
	}
}

func Reduce(acc stream.T, fn stream.ReduceFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			acc = fn(acc, data)
			return nil
		},
		OnCompleted: func(emitter stream.Emitter) {
			emitter.Emit(acc)
		},
	}
}

func Flatten() stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			dv := reflect.ValueOf(data)
			if dv.Kind() == reflect.Slice || dv.Kind() == reflect.Ptr && dv.Elem().Kind() == reflect.Slice {
				for i := 0; i < dv.Len(); i++ {
					emitter.Emit(dv.Index(i).Interface())
				}
			} else {
				emitter.Emit(data)
			}
			return nil
		},
	}
}

func Batch(size int) stream.Transformer {
	return BatchBy(&batch{size: size})
}

func BatchBy(batch stream.Batch) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			batch.Add(data)
			if batch.Full() {
				batch.Commit(emitter)
			}
			return nil
		},
		OnCompleted: func(emitter stream.Emitter) {
			if !batch.Empty() {
				batch.Commit(emitter)
			}
		},
	}
}

func Each(fn stream.EachFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, emitter stream.Emitter) error {
			fn(data)
			emitter.Emit(data)
			return nil
		},
	}
}
