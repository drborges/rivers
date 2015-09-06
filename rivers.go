package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
)

type Pipeline struct {
	Context stream.Context
	Stream  stream.Readable
}

func From(producer stream.Producer) *Pipeline {
	context := NewContext()

	if bindable, ok := producer.(stream.Bindable); ok {
		bindable.Bind(context)
	}

	return &Pipeline{
		Context: context,
		Stream:  producer.Produce(),
	}
}

func FromStream(readable stream.Readable) *Pipeline {
	return &Pipeline{
		Context: NewContext(),
		Stream:  readable,
	}
}

func FromRange(from, to int) *Pipeline {
	return From(producers.FromRange(from, to))
}

func FromData(data ...stream.T) *Pipeline {
	return From(producers.FromData(data...))
}

func FromSlice(slice stream.T) *Pipeline {
	return From(producers.FromSlice(slice))
}

func (s *Pipeline) Split() (*Pipeline, *Pipeline) {
	streams := s.SplitN(2)
	return streams[0], streams[1]
}

func (s *Pipeline) SplitN(n int) []*Pipeline {
	pipelines := make([]*Pipeline, n)
	writables := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		readable, writable := stream.New(cap(s.Stream))
		pipelines[i] = &Pipeline{s.Context, readable}
		writables[i] = writable
	}
	dispatchers.New(s.Context).Always().Dispatch(s.Stream, writables...)
	return pipelines
}

func (s *Pipeline) Partition(fn stream.PredicateFn) (*Pipeline, *Pipeline) {
	lhsIn, lhsOut := stream.New(cap(s.Stream))
	rhsIn := dispatchers.New(s.Context).If(fn).Dispatch(s.Stream, lhsOut)

	return &Pipeline{s.Context, lhsIn}, &Pipeline{s.Context, rhsIn}
}

func (s *Pipeline) Merge(readables ...stream.Readable) *Pipeline {
	combiner := combiners.FIFO()
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeMerged := []stream.Readable{s.Stream}
	toBeMerged = append(toBeMerged, readables...)
	return &Pipeline{
		Context: s.Context,
		Stream: combiner.Combine(toBeMerged...),
	}
}

func (s *Pipeline) Zip(readables ...stream.Readable) *Pipeline {
	combiner := combiners.Zip()
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeZipped := []stream.Readable{s.Stream}
	toBeZipped = append(toBeZipped, readables...)
	return &Pipeline{
		Context: s.Context,
		Stream: combiner.Combine(toBeZipped...),
	}
}

func (s *Pipeline) ZipBy(fn stream.ReduceFn, readables ...stream.Readable) *Pipeline {
	combiner := combiners.ZipBy(fn)
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeZipped := []stream.Readable{s.Stream}
	toBeZipped = append(toBeZipped, readables...)
	return &Pipeline{
		Context: s.Context,
		Stream:  combiner.Combine(toBeZipped...),
	}
}

func (s *Pipeline) Dispatch(writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context: s.Context,
		Stream: dispatchers.New(s.Context).Always().Dispatch(s.Stream, writables...),
	}
}

func (s *Pipeline) DispatchIf(fn stream.PredicateFn, writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context: s.Context,
		Stream: dispatchers.New(s.Context).If(fn).Dispatch(s.Stream, writables...),
	}
}

func (s *Pipeline) Apply(transformer stream.Transformer) *Pipeline {
	if bindable, ok := transformer.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	return &Pipeline{
		Stream:  transformer.Transform(s.Stream),
		Context: s.Context,
	}
}

func (s *Pipeline) Filter(fn stream.PredicateFn) *Pipeline {
	return s.Apply(transformers.Filter(fn))
}

func (s *Pipeline) OnData(fn stream.OnDataFn) *Pipeline {
	return s.Apply(transformers.OnData(fn))
}

func (s *Pipeline) Map(fn stream.MapFn) *Pipeline {
	return s.Apply(transformers.Map(fn))
}

func (s *Pipeline) Each(fn stream.EachFn) *Pipeline {
	return s.Apply(transformers.Each(fn))
}

func (s *Pipeline) FindBy(fn stream.PredicateFn) *Pipeline {
	return s.Apply(transformers.FindBy(fn))
}

func (s *Pipeline) TakeFirst(n int) *Pipeline {
	return s.Apply(transformers.TakeFirst(n))
}

func (s *Pipeline) Take(fn stream.PredicateFn) *Pipeline {
	return s.Apply(transformers.Take(fn))
}

func (s *Pipeline) DropFirst(n int) *Pipeline {
	return s.Apply(transformers.DropFirst(n))
}

func (s *Pipeline) Drop(fn stream.PredicateFn) *Pipeline {
	return s.Apply(transformers.Drop(fn))
}

func (s *Pipeline) Reduce(acc stream.T, fn stream.ReduceFn) *Pipeline {
	return s.Apply(transformers.Reduce(acc, fn))
}

func (s *Pipeline) Flatten() *Pipeline {
	return s.Apply(transformers.Flatten())
}

func (s *Pipeline) SortBy(fn stream.SortByFn) *Pipeline {
	return s.Apply(transformers.SortBy(fn))
}

func (s *Pipeline) Batch(size int) *Pipeline {
	return s.Apply(transformers.Batch(size))
}

func (s *Pipeline) BatchBy(batch stream.Batch) *Pipeline {
	return s.Apply(transformers.BatchBy(batch))
}

func (s *Pipeline) Then(consumer stream.Consumer) error {
	if bindable, ok := consumer.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}
	consumer.Consume(s.Stream)
	return s.Context.Err()
}

func (s *Pipeline) Collect() ([]stream.T, error) {
	var data []stream.T
	return data, s.CollectAs(&data)
}

func (s *Pipeline) CollectAs(data interface{}) error {
	return s.Then(consumers.ItemsCollector(data))
}

func (s *Pipeline) CollectFirst() (stream.T, error) {
	var data stream.T
	return data, s.CollectFirstAs(&data)
}

func (s *Pipeline) CollectFirstAs(data interface{}) error {
	return s.TakeFirst(1).Then(consumers.LastItemCollector(data))
}

func (s *Pipeline) CollectLast() (stream.T, error) {
	var data stream.T
	return data, s.CollectLastAs(&data)
}

func (s *Pipeline) CollectLastAs(data interface{}) error {
	return s.Then(consumers.LastItemCollector(data))
}

func (s *Pipeline) CollectBy(fn stream.EachFn) error {
	return s.Then(consumers.CollectBy(fn))
}

func (s *Pipeline) Drain() error {
	return s.Then(consumers.Drainer())
}
