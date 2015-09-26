package stream

import (
	"errors"
	"time"
)

var (
	Done    = errors.New("Context is done")
	Timeout = errors.New("Context has timed out")
)

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
	Deadline() time.Duration
	SetDeadline(time.Duration)
	Failure() <-chan struct{}
	Done() <-chan struct{}
}

// a.k.a Source
type Producer interface {
	Attachable
	Produce() (out Readable)
}

// a.k.a Sink
type Consumer interface {
	Attachable
	Consume(in Readable)
}

type Transformer interface {
	Attachable
	Transform(in Readable) (out Readable)
}

// FIFO, Zip, InOrder combiner
type Combiner interface {
	Attachable
	Combine(in ...Readable) (out Readable)
}

type Dispatcher interface {
	Attachable
	Dispatch(from Readable, to ...Writable) (out Readable)
}

type Emitter interface {
	Emit(data T)
}

type Attachable interface {
	Attach(Context)
}

type Batch interface {
	Commit(Emitter)
	Full() bool
	Empty() bool
	Add(data T)
}

type Groups map[T][]T

func (groups Groups) Empty() bool {
	return len(groups) == 0
}

func (groups Groups) HasGroup(key T) bool {
	_, contains := groups[key]
	return contains
}

func (groups Groups) HasItem(val T) bool {
	for _, group := range groups {
		for _, item := range group {
			if item == val {
				return true
			}
		}
	}
	return false
}
