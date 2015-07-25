package rx

type T interface{}
type InStream <-chan T
type OutStream chan<- T
type MapFn func(T) T
type HandleFn func(T)
type PredicateFn func(T) bool
type SortByFn func(a, b T) bool
type ReduceFn func(acc, next T) (result T)

type Context interface {
	Close()
	Recover()
	Closed() <-chan error
}

// a.k.a Source
type Producer interface {
	Produce() (out InStream)
}

// a.k.a Sink
// TODO better evaluate the need for this interface
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
