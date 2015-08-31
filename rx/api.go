package rx

type T interface{}
type Readable <-chan T
type Writable chan<- T
type MapFn func(T) T
type EachFn func(T)
type PredicateFn func(T) bool
type SortByFn func(a, b T) bool
type OnDataFn func(data T, out Writable)
type ReduceFn func(acc, next T) (result T)

type Context interface {
	Close()
	Recover()
	Err() error
	Closed() <-chan error
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
	Dispatch(from Readable, to ...Writable) (sink Readable)
}

type Batch interface {
	Commit(Writable)
	Full() bool
	Empty() bool
	Add(data T)
}
