package stream

import "errors"

var Done = errors.New("Context closed")

type T interface{}
type Readable <-chan T
type Writable chan<- T
type MapFn func(T) T
type EachFn func(T)
type PredicateFn func(T) bool
type SortByFn func(a, b T) bool
type OnDataFn func(data T, emitter Emitter)
type ReduceFn func(acc, next T) (result T)

type Context interface {
	Close(err error)
	Recover()
	Err() error
	Done() <-chan struct{}
}

// a.k.a Source
type Producer interface {
	Produce() (out Readable)
}

// a.k.a Sink
type Consumer interface {
	Consume(in Readable)
}

type Transformer interface {
	Transform(in Readable) (out Readable)
}

// FIFO, Zip, InOrder combiner
type Combiner interface {
	Combine(in ...Readable) (out Readable)
}

type Dispatcher interface {
	Dispatch(from Readable, to ...Writable) (out Readable)
}

type Emitter interface {
	Emit(data T)
}

type Bindable interface {
	Bind(Context)
}

type Batch interface {
	Commit(Emitter)
	Full() bool
	Empty() bool
	Add(data T)
}
