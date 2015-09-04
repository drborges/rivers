package transformers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

func Filter(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			if fn(data) {
				w <- data
			}
			return nil
		},
	}
}

func FindBy(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			if fn(data) {
				w <- data
				return stream.Done
			}
			return nil
		},
	}
}

func TakeFirst(n int) stream.Transformer {
	taken := 0
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			if taken >= n {
				return stream.Done
			}

			w <- data
			taken++
			return nil
		},
	}
}

func Take(fn stream.PredicateFn) stream.Transformer {
	return Filter(fn)
}

func Drop(fn stream.PredicateFn) stream.Transformer {
	return Filter(func(data stream.T) bool { return !fn(data) })
}

func Map(fn stream.MapFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			w <- fn(data)
			return nil
		},
	}
}

func OnData(fn stream.OnDataFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			fn(data, w)
			return nil
		},
	}
}

func Reduce(acc stream.T, fn stream.ReduceFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			acc = fn(acc, data)
			return nil
		},
		OnCompleted: func(w stream.Writable) {
			w <- acc
		},
	}
}

func Flatten() stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			dv := reflect.ValueOf(data)
			if dv.Kind() == reflect.Slice || dv.Kind() == reflect.Ptr && dv.Elem().Kind() == reflect.Slice {
				for i := 0; i < dv.Len(); i++ {
					w <- dv.Index(i).Interface()
				}
			} else {
				w <- data
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
		OnNext: func(data stream.T, w stream.Writable) error {
			batch.Add(data)
			if batch.Full() {
				batch.Commit(w)
			}
			return nil
		},
		OnCompleted: func(w stream.Writable) {
			if !batch.Empty() {
				batch.Commit(w)
			}
		},
	}
}

func SortBy(sorter stream.SortByFn) stream.Transformer {
	items := []stream.T{}
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			items = append(items, data)
			return nil
		},
		OnCompleted: func(w stream.Writable) {
			sorter.Sort(items)
			for _, item := range items {
				w <- item
			}
		},
	}
}

func Each(fn stream.EachFn) stream.Transformer {
	return &Observer{
		OnNext: func(data stream.T, w stream.Writable) error {
			fn(data)
			w <- data
			return nil
		},
	}
}
