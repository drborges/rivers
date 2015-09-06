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

func (pipeline *Pipeline) Split() (*Pipeline, *Pipeline) {
	pipelines := pipeline.SplitN(2)
	return pipelines[0], pipelines[1]
}

func (pipeline *Pipeline) SplitN(n int) []*Pipeline {
	pipelines := make([]*Pipeline, n)
	writables := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		readable, writable := stream.New(cap(pipeline.Stream))
		pipelines[i] = &Pipeline{pipeline.Context, readable}
		writables[i] = writable
	}
	dispatchers.New(pipeline.Context).Always().Dispatch(pipeline.Stream, writables...)
	return pipelines
}

func (pipeline *Pipeline) Partition(fn stream.PredicateFn) (*Pipeline, *Pipeline) {
	lhsIn, lhsOut := stream.New(cap(pipeline.Stream))
	rhsIn := dispatchers.New(pipeline.Context).If(fn).Dispatch(pipeline.Stream, lhsOut)

	return &Pipeline{pipeline.Context, lhsIn}, &Pipeline{pipeline.Context, rhsIn}
}

func (pipeline *Pipeline) Dispatch(writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context: pipeline.Context,
		Stream:  dispatchers.New(pipeline.Context).Always().Dispatch(pipeline.Stream, writables...),
	}
}

func (pipeline *Pipeline) DispatchIf(fn stream.PredicateFn, writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context: pipeline.Context,
		Stream:  dispatchers.New(pipeline.Context).If(fn).Dispatch(pipeline.Stream, writables...),
	}
}

func (pipeline *Pipeline) Merge(pipelines ...*Pipeline) *Pipeline {
	return pipeline.Combine(combiners.FIFO(), pipelines)
}

func (pipeline *Pipeline) Zip(pipelines ...*Pipeline) *Pipeline {
	return pipeline.Combine(combiners.Zip(), pipelines)
}

func (pipeline *Pipeline) ZipBy(fn stream.ReduceFn, pipelines ...*Pipeline) *Pipeline {
	return pipeline.Combine(combiners.ZipBy(fn), pipelines)
}

func (pipeline *Pipeline) Combine(combiner stream.Combiner, pipelines []*Pipeline) *Pipeline {
	if bindable, ok := combiner.(stream.Bindable); ok {
		bindable.Bind(pipeline.Context)
	}

	readables := []stream.Readable{pipeline.Stream}
	for _, p := range pipelines {
		readables = append(readables, p.Stream)
	}

	return &Pipeline{
		Context: pipeline.Context,
		Stream:  combiner.Combine(readables...),
	}
}

func (pipeline *Pipeline) Apply(transformer stream.Transformer) *Pipeline {
	if bindable, ok := transformer.(stream.Bindable); ok {
		bindable.Bind(pipeline.Context)
	}

	return &Pipeline{
		Stream:  transformer.Transform(pipeline.Stream),
		Context: pipeline.Context,
	}
}

func (pipeline *Pipeline) Filter(fn stream.PredicateFn) *Pipeline {
	return pipeline.Apply(transformers.Filter(fn))
}

func (pipeline *Pipeline) OnData(fn stream.OnDataFn) *Pipeline {
	return pipeline.Apply(transformers.OnData(fn))
}

func (pipeline *Pipeline) Map(fn stream.MapFn) *Pipeline {
	return pipeline.Apply(transformers.Map(fn))
}

func (pipeline *Pipeline) Each(fn stream.EachFn) *Pipeline {
	return pipeline.Apply(transformers.Each(fn))
}

func (pipeline *Pipeline) FindBy(fn stream.PredicateFn) *Pipeline {
	return pipeline.Apply(transformers.FindBy(fn))
}

func (pipeline *Pipeline) TakeFirst(n int) *Pipeline {
	return pipeline.Apply(transformers.TakeFirst(n))
}

func (pipeline *Pipeline) Take(fn stream.PredicateFn) *Pipeline {
	return pipeline.Apply(transformers.Take(fn))
}

func (pipeline *Pipeline) DropFirst(n int) *Pipeline {
	return pipeline.Apply(transformers.DropFirst(n))
}

func (pipeline *Pipeline) Drop(fn stream.PredicateFn) *Pipeline {
	return pipeline.Apply(transformers.Drop(fn))
}

func (pipeline *Pipeline) Reduce(acc stream.T, fn stream.ReduceFn) *Pipeline {
	return pipeline.Apply(transformers.Reduce(acc, fn))
}

func (pipeline *Pipeline) Flatten() *Pipeline {
	return pipeline.Apply(transformers.Flatten())
}

func (pipeline *Pipeline) SortBy(fn stream.SortByFn) *Pipeline {
	return pipeline.Apply(transformers.SortBy(fn))
}

func (pipeline *Pipeline) Batch(size int) *Pipeline {
	return pipeline.Apply(transformers.Batch(size))
}

func (pipeline *Pipeline) BatchBy(batch stream.Batch) *Pipeline {
	return pipeline.Apply(transformers.BatchBy(batch))
}

func (pipeline *Pipeline) Then(consumer stream.Consumer) error {
	if bindable, ok := consumer.(stream.Bindable); ok {
		bindable.Bind(pipeline.Context)
	}
	consumer.Consume(pipeline.Stream)
	return pipeline.Context.Err()
}

func (pipeline *Pipeline) Collect() ([]stream.T, error) {
	var data []stream.T
	return data, pipeline.CollectAs(&data)
}

func (pipeline *Pipeline) CollectAs(data interface{}) error {
	return pipeline.Then(consumers.ItemsCollector(data))
}

func (pipeline *Pipeline) CollectFirst() (stream.T, error) {
	var data stream.T
	return data, pipeline.CollectFirstAs(&data)
}

func (pipeline *Pipeline) CollectFirstAs(data interface{}) error {
	return pipeline.TakeFirst(1).Then(consumers.LastItemCollector(data))
}

func (pipeline *Pipeline) CollectLast() (stream.T, error) {
	var data stream.T
	return data, pipeline.CollectLastAs(&data)
}

func (pipeline *Pipeline) CollectLastAs(data interface{}) error {
	return pipeline.Then(consumers.LastItemCollector(data))
}

func (pipeline *Pipeline) CollectBy(fn stream.EachFn) error {
	return pipeline.Then(consumers.CollectBy(fn))
}

func (pipeline *Pipeline) Drain() error {
	return pipeline.Then(consumers.Drainer())
}
