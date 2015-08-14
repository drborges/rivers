package rx

type T interface{}
type InStream <-chan T
type OutStream chan<- T
type MapFn func(T) T
type HandleFn func(T)
type PredicateFn func(T) bool
type SortByFn func(a, b T) bool
type OnDataFn func(data T, out OutStream)
type ReduceFn func(acc, next T) (result T)

type Context interface {
	Close()
	Recover()
	Err() error
	Closed() <-chan error
}

// a.k.a Source
type Producer interface {
	Produce() (out InStream)
}

// a.k.a Sink
type Consumer interface {
	Consume(in InStream)
}

type Transformer interface {
	Transform(in InStream) (out InStream)
}

// FIFO, Zip, InOrder combiner
type Combiner interface {
	Combine(in ...InStream) (out InStream)
}

type Dispatcher interface {
	Dispatch(from InStream, to ...OutStream) (sink InStream)
}

type Batch interface {
	Commit(OutStream)
	Full() bool
	Empty() bool
	Add(data T)
}
