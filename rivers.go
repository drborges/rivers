package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
)

type Stream struct {
	readable     stream.Readable
	Context      stream.Context
	consumers    *consumers.Builder
	combiners    *combiners.Builder
	dispatchers  *dispatchers.Builder
	transformers *transformers.Builder
}

func From(producer stream.Producer) *Stream {
	context := NewContext()

	if bindable, ok := producer.(stream.Bindable); ok {
		bindable.Bind(context)
	}

	return NewWith(context).newFrom(producer.Produce())
}

func FromStream(readable stream.Readable) *Stream {
	return NewWith(NewContext()).newFrom(readable)
}

func FromRange(from, to int) *Stream {
	return From(producers.FromRange(from, to))
}

func FromData(data ...stream.T) *Stream {
	return From(producers.FromData(data...))
}

func FromSlice(slice stream.T) *Stream {
	return From(producers.FromSlice(slice))
}

func Merge(streams ...*Stream) *Stream {
	readables := []stream.Readable{}
	for _, s := range streams {
		readables = append(readables, s.Sink())
	}
	return NewWith(NewContext()).Merge(readables...)
}

func Zip(streams ...*Stream) *Stream {
	readables := []stream.Readable{}
	for _, s := range streams {
		readables = append(readables, s.Sink())
	}
	return NewWith(NewContext()).Zip(readables...)
}

func ZipBy(fn stream.ReduceFn, streams ...*Stream) *Stream {
	readables := []stream.Readable{}
	for _, s := range streams {
		readables = append(readables, s.Sink())
	}
	return NewWith(NewContext()).ZipBy(fn, readables...)
}

func (s *Stream) newFrom(readable stream.Readable) *Stream {
	return &Stream{
		readable:     readable,
		Context:      s.Context,
		consumers:    s.consumers,
		combiners:    s.combiners,
		dispatchers:  s.dispatchers,
		transformers: s.transformers,
	}
}

func NewWith(context stream.Context) *Stream {
	return &Stream{
		Context:      context,
		consumers:    consumers.New(context),
		combiners:    combiners.New(context),
		dispatchers:  dispatchers.New(context),
		transformers: transformers.New(context),
	}
}

func (s *Stream) Split() (*Stream, *Stream) {
	streams := s.SplitN(2)
	return streams[0], streams[1]
}

func (s *Stream) SplitN(n int) []*Stream {
	streams := make([]*Stream, n)
	writables := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		readable, writable := stream.New(cap(s.readable))
		streams[i] = s.newFrom(readable)
		writables[i] = writable
	}
	s.dispatchers.Always().Dispatch(s.readable, writables...)
	return streams
}

func (s *Stream) Partition(fn stream.PredicateFn) (*Stream, *Stream) {
	lhsIn, lhsOut := stream.New(cap(s.readable))
	rhsIn := s.dispatchers.If(fn).Dispatch(s.readable, lhsOut)

	return s.newFrom(lhsIn), s.newFrom(rhsIn)
}

func (s *Stream) Merge(readables ...stream.Readable) *Stream {
	return s.newFrom(s.combiners.FIFO().Combine(readables...))
}

func (s *Stream) Zip(readables ...stream.Readable) *Stream {
	return s.newFrom(s.combiners.Zip().Combine(readables...))
}

func (s *Stream) ZipBy(fn stream.ReduceFn, readables ...stream.Readable) *Stream {
	return s.newFrom(s.combiners.ZipBy(fn).Combine(readables...))
}

func (s *Stream) Dispatch(writables ...stream.Writable) *Stream {
	return s.newFrom(s.dispatchers.Always().Dispatch(s.readable, writables...))
}

func (s *Stream) DispatchIf(fn stream.PredicateFn, writables ...stream.Writable) *Stream {
	return s.newFrom(s.dispatchers.If(fn).Dispatch(s.readable, writables...))
}

func (s *Stream) Apply(t stream.Transformer) *Stream {
	return s.newFrom(t.Transform(s.readable))
}

func (s *Stream) Filter(fn stream.PredicateFn) *Stream {
	return s.Apply(s.transformers.Filter(fn))
}

func (s *Stream) OnData(fn stream.OnDataFn) *Stream {
	return s.Apply(s.transformers.OnData(fn))
}

func (s *Stream) Map(fn stream.MapFn) *Stream {
	return s.Apply(s.transformers.Map(fn))
}

func (s *Stream) Each(fn stream.EachFn) *Stream {
	return s.Apply(s.transformers.Each(fn))
}

func (s *Stream) FindBy(fn stream.PredicateFn) *Stream {
	return s.Apply(s.transformers.FindBy(fn))
}

func (s *Stream) TakeFirst(n int) *Stream {
	return s.Apply(s.transformers.TakeFirst(n))
}

func (s *Stream) Take(fn stream.PredicateFn) *Stream {
	return s.Apply(s.transformers.Take(fn))
}

func (s *Stream) Drop(fn stream.PredicateFn) *Stream {
	return s.Apply(s.transformers.Drop(fn))
}

func (s *Stream) Reduce(acc stream.T, fn stream.ReduceFn) *Stream {
	return s.Apply(s.transformers.Reduce(acc, fn))
}

func (s *Stream) Flatten() *Stream {
	return s.Apply(s.transformers.Flatten())
}

func (s *Stream) SortBy(fn stream.SortByFn) *Stream {
	return s.Apply(s.transformers.SortBy(fn))
}

func (s *Stream) Batch(size int) *Stream {
	return s.Apply(s.transformers.Batch(size))
}

func (s *Stream) BatchBy(batch stream.Batch) *Stream {
	return s.Apply(s.transformers.BatchBy(batch))
}

func (s *Stream) Sink() stream.Readable {
	return s.readable
}

func (s *Stream) Collect() ([]stream.T, error) {
	var data []stream.T
	return data, s.CollectAs(&data)
}

func (s *Stream) CollectAs(data interface{}) error {
	s.consumers.ItemsCollector(data).Consume(s.readable)
	return s.Context.Err()
}

func (s *Stream) CollectFirst() (stream.T, error) {
	var data stream.T
	return data, s.CollectFirstAs(&data)
}

func (s *Stream) CollectFirstAs(data interface{}) error {
	s.consumers.LastItemCollector(data).Consume(s.TakeFirst(1).Sink())
	return s.Context.Err()
}

func (s *Stream) CollectLast() (stream.T, error) {
	var data stream.T
	return data, s.CollectLastAs(&data)
}

func (s *Stream) CollectLastAs(data interface{}) error {
	s.consumers.LastItemCollector(data).Consume(s.readable)
	return s.Context.Err()
}

func (s *Stream) Drain() error {
	s.consumers.Drainer().Consume(s.readable)
	return s.Context.Err()
}
