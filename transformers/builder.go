package transformers

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Filter(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			if fn(data) {
				w <- data
			}
			return nil
		},
	}
}

func (b *Builder) FindBy(fn stream.PredicateFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			if fn(data) {
				w <- data
				return stream.Done
			}
			return nil
		},
	}
}

func (b *Builder) Take(n int) stream.Transformer {
	taken := 0
	return &Observer{
		Context: b.context,
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

func (b *Builder) TakeIf(fn stream.PredicateFn) stream.Transformer {
	return b.Filter(fn)
}

func (b *Builder) DropIf(fn stream.PredicateFn) stream.Transformer {
	return b.Filter(func(data stream.T) bool { return !fn(data) })
}

func (b *Builder) Map(fn stream.MapFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			w <- fn(data)
			return nil
		},
	}
}

func (b *Builder) OnData(fn stream.OnDataFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			fn(data, w)
			return nil
		},
	}
}

func (b *Builder) Reduce(acc stream.T, fn stream.ReduceFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			acc = fn(acc, data)
			return nil
		},
		OnCompleted: func(w stream.Writable) {
			w <- acc
		},
	}
}

func (b *Builder) Flatten() stream.Transformer {
	return &Observer{
		Context: b.context,
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

func (b *Builder) Batch(size int) stream.Transformer {
	return b.BatchBy(&batch{size: size})
}

func (b *Builder) BatchBy(batch stream.Batch) stream.Transformer {
	return &Observer{
		Context: b.context,
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

func (b *Builder) SortBy(sorter stream.SortByFn) stream.Transformer {
	items := []stream.T{}
	return &Observer{
		Context: b.context,
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

func (b *Builder) Each(fn stream.EachFn) stream.Transformer {
	return &Observer{
		Context: b.context,
		OnNext: func(data stream.T, w stream.Writable) error {
			fn(data)
			w <- data
			return nil
		},
	}
}
