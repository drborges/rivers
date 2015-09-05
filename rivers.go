package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
)

type Stage struct {
	Context  stream.Context
	Stream stream.Readable
}

func From(producer stream.Producer) *Stage {
	context := NewContext()

	if bindable, ok := producer.(stream.Bindable); ok {
		bindable.Bind(context)
	}

	return NewStage(context).newStage(producer.Produce())
}

func FromStream(readable stream.Readable) *Stage {
	return NewStage(NewContext()).newStage(readable)
}

func FromRange(from, to int) *Stage {
	return From(producers.FromRange(from, to))
}

func FromData(data ...stream.T) *Stage {
	return From(producers.FromData(data...))
}

func FromSlice(slice stream.T) *Stage {
	return From(producers.FromSlice(slice))
}

func (s *Stage) newStage(readable stream.Readable) *Stage {
	return &Stage{
		Stream: readable,
		Context:  s.Context,
	}
}

func NewStage(context stream.Context) *Stage {
	return &Stage{
		Context:  context,
		Stream: stream.NewEmpty(),
	}
}

func (s *Stage) Split() (*Stage, *Stage) {
	streams := s.SplitN(2)
	return streams[0], streams[1]
}

func (s *Stage) SplitN(n int) []*Stage {
	streams := make([]*Stage, n)
	writables := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		readable, writable := stream.New(cap(s.Stream))
		streams[i] = s.newStage(readable)
		writables[i] = writable
	}
	dispatchers.New(s.Context).Always().Dispatch(s.Stream, writables...)
	return streams
}

func (s *Stage) Partition(fn stream.PredicateFn) (*Stage, *Stage) {
	lhsIn, lhsOut := stream.New(cap(s.Stream))
	rhsIn := dispatchers.New(s.Context).If(fn).Dispatch(s.Stream, lhsOut)

	return s.newStage(lhsIn), s.newStage(rhsIn)
}

func (s *Stage) Merge(readables ...stream.Readable) *Stage {
	combiner := combiners.FIFO()
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeMerged := []stream.Readable{s.Stream}
	toBeMerged = append(toBeMerged, readables...)
	return s.newStage(combiner.Combine(toBeMerged...))
}

func (s *Stage) Zip(readables ...stream.Readable) *Stage {
	combiner := combiners.Zip()
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeZipped := []stream.Readable{s.Stream}
	toBeZipped = append(toBeZipped, readables...)
	return s.newStage(combiner.Combine(toBeZipped...))
}

func (s *Stage) ZipBy(fn stream.ReduceFn, readables ...stream.Readable) *Stage {
	combiner := combiners.ZipBy(fn)
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	toBeZipped := []stream.Readable{s.Stream}
	toBeZipped = append(toBeZipped, readables...)
	return s.newStage(combiner.Combine(toBeZipped...))
}

func (s *Stage) Dispatch(writables ...stream.Writable) *Stage {
	return s.newStage(dispatchers.New(s.Context).Always().Dispatch(s.Stream, writables...))
}

func (s *Stage) DispatchIf(fn stream.PredicateFn, writables ...stream.Writable) *Stage {
	return s.newStage(dispatchers.New(s.Context).If(fn).Dispatch(s.Stream, writables...))
}

func (s *Stage) Apply(transformer stream.Transformer) *Stage {
	if bindable, ok := transformer.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}

	return &Stage{
		Stream: transformer.Transform(s.Stream),
		Context:  s.Context,
	}
}

func (s *Stage) Filter(fn stream.PredicateFn) *Stage {
	return s.Apply(transformers.Filter(fn))
}

func (s *Stage) OnData(fn stream.OnDataFn) *Stage {
	return s.Apply(transformers.OnData(fn))
}

func (s *Stage) Map(fn stream.MapFn) *Stage {
	return s.Apply(transformers.Map(fn))
}

func (s *Stage) Each(fn stream.EachFn) *Stage {
	return s.Apply(transformers.Each(fn))
}

func (s *Stage) FindBy(fn stream.PredicateFn) *Stage {
	return s.Apply(transformers.FindBy(fn))
}

func (s *Stage) TakeFirst(n int) *Stage {
	return s.Apply(transformers.TakeFirst(n))
}

func (s *Stage) Take(fn stream.PredicateFn) *Stage {
	return s.Apply(transformers.Take(fn))
}

func (s *Stage) DropFirst(n int) *Stage {
	return s.Apply(transformers.DropFirst(n))
}

func (s *Stage) Drop(fn stream.PredicateFn) *Stage {
	return s.Apply(transformers.Drop(fn))
}

func (s *Stage) Reduce(acc stream.T, fn stream.ReduceFn) *Stage {
	return s.Apply(transformers.Reduce(acc, fn))
}

func (s *Stage) Flatten() *Stage {
	return s.Apply(transformers.Flatten())
}

func (s *Stage) SortBy(fn stream.SortByFn) *Stage {
	return s.Apply(transformers.SortBy(fn))
}

func (s *Stage) Batch(size int) *Stage {
	return s.Apply(transformers.Batch(size))
}

func (s *Stage) BatchBy(batch stream.Batch) *Stage {
	return s.Apply(transformers.BatchBy(batch))
}

func (s *Stage) Then(consumer stream.Consumer) error {
	if bindable, ok := consumer.(stream.Bindable); ok {
		bindable.Bind(s.Context)
	}
	consumer.Consume(s.Stream)
	return s.Context.Err()
}

func (s *Stage) Collect() ([]stream.T, error) {
	var data []stream.T
	return data, s.CollectAs(&data)
}

func (s *Stage) CollectAs(data interface{}) error {
	return s.Then(consumers.ItemsCollector(data))
}

func (s *Stage) CollectFirst() (stream.T, error) {
	var data stream.T
	return data, s.CollectFirstAs(&data)
}

func (s *Stage) CollectFirstAs(data interface{}) error {
	return s.TakeFirst(1).Then(consumers.LastItemCollector(data))
}

func (s *Stage) CollectLast() (stream.T, error) {
	var data stream.T
	return data, s.CollectLastAs(&data)
}

func (s *Stage) CollectLastAs(data interface{}) error {
	return s.Then(consumers.LastItemCollector(data))
}

func (s *Stage) CollectBy(fn stream.EachFn) error {
	return s.Then(consumers.CollectBy(fn))
}

func (s *Stage) Drain() error {
	return s.Then(consumers.Drainer())
}
