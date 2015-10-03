package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	"io"
	"time"
)

type Pipeline struct {
	Context  stream.Context
	Stream   stream.Readable
	parallel bool
}

func From(producer stream.Producer) *Pipeline {
	context := NewContext()
	producer.Attach(context)

	return &Pipeline{
		Context: context,
		Stream:  producer.Produce(),
	}
}

func FromRange(from, to int) *Pipeline {
	return From(producers.FromRange(from, to))
}

func FromReader(r io.Reader) *Pipeline {
	return From(producers.FromReader(r))
}

func FromData(data ...stream.T) *Pipeline {
	return From(producers.FromData(data...))
}

func FromSlice(slice stream.T) *Pipeline {
	return From(producers.FromSlice(slice))
}

func (pipeline *Pipeline) Parallel() *Pipeline {
	pipeline.parallel = true
	return pipeline
}

func (pipeline *Pipeline) Deadline(duration time.Duration) *Pipeline {
	pipeline.Context.SetDeadline(duration)
	return pipeline
}

func (pipeline *Pipeline) Split() (*Pipeline, *Pipeline) {
	pipelines := pipeline.SplitN(2)
	return pipelines[0], pipelines[1]
}

func (pipeline *Pipeline) SplitN(n int) []*Pipeline {
	pipelines := make([]*Pipeline, n)
	writables := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		readable, writable := stream.New(pipeline.Stream.Capacity())
		writables[i] = writable
		pipelines[i] = &Pipeline{
			Context:  pipeline.Context,
			Stream:   readable,
			parallel: pipeline.parallel,
		}
	}
	dispatchers.New(pipeline.Context).Always().Dispatch(pipeline.Stream, writables...)
	return pipelines
}

func (pipeline *Pipeline) Partition(fn stream.PredicateFn) (*Pipeline, *Pipeline) {
	lhsIn, lhsOut := stream.New(pipeline.Stream.Capacity())
	rhsIn := dispatchers.New(pipeline.Context).If(fn).Dispatch(pipeline.Stream, lhsOut)
	lhsPipeline := &Pipeline{Context: pipeline.Context, Stream: lhsIn, parallel: pipeline.parallel}
	rhsPipeline := &Pipeline{Context: pipeline.Context, Stream: rhsIn, parallel: pipeline.parallel}
	return lhsPipeline, rhsPipeline
}

func (pipeline *Pipeline) Dispatch(writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context:  pipeline.Context,
		Stream:   dispatchers.New(pipeline.Context).Always().Dispatch(pipeline.Stream, writables...),
		parallel: pipeline.parallel,
	}
}

func (pipeline *Pipeline) DispatchIf(fn stream.PredicateFn, writables ...stream.Writable) *Pipeline {
	return &Pipeline{
		Context:  pipeline.Context,
		Stream:   dispatchers.New(pipeline.Context).If(fn).Dispatch(pipeline.Stream, writables...),
		parallel: pipeline.parallel,
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
	combiner.Attach(pipeline.Context)

	readables := []stream.Readable{pipeline.Stream}
	for _, p := range pipelines {
		readables = append(readables, p.Stream)
	}

	return &Pipeline{
		Context:  pipeline.Context,
		Stream:   combiner.Combine(readables...),
		parallel: pipeline.parallel,
	}
}

func (pipeline *Pipeline) Apply(transformer stream.Transformer) *Pipeline {
	transformer.Attach(pipeline.Context)

	return &Pipeline{
		Stream:   transformer.Transform(pipeline.Stream),
		Context:  pipeline.Context,
		parallel: pipeline.parallel,
	}
}

func (pipeline *Pipeline) ApplyParallel(transformer stream.Transformer) *Pipeline {
	parallelPipelineCount := 0
	if pipeline.parallel {
		parallelPipelineCount = pipeline.Stream.Capacity() - 1
	}

	parallelPipelines := make([]*Pipeline, parallelPipelineCount)
	for i, _ := range parallelPipelines {
		parallelPipelines[i] = pipeline.Apply(transformer)
	}

	return pipeline.Apply(transformer).Merge(parallelPipelines...)
}

func (pipeline *Pipeline) Filter(fn stream.PredicateFn) *Pipeline {
	return pipeline.ApplyParallel(transformers.Filter(fn))
}

func (pipeline *Pipeline) OnData(fn stream.OnDataFn) *Pipeline {
	return pipeline.ApplyParallel(transformers.OnData(fn))
}

func (pipeline *Pipeline) Map(fn stream.MapFn) *Pipeline {
	return pipeline.ApplyParallel(transformers.Map(fn))
}

func (pipeline *Pipeline) FlatMap(fn stream.MapFn) *Pipeline {
	return pipeline.ApplyParallel(transformers.Map(fn)).Flatten()
}

func (pipeline *Pipeline) Each(fn stream.EachFn) *Pipeline {
	return pipeline.ApplyParallel(transformers.Each(fn))
}

func (pipeline *Pipeline) Find(subject stream.T) *Pipeline {
	return pipeline.Apply(transformers.FindBy(func(data stream.T) bool {
		return data == subject
	}))
}

func (pipeline *Pipeline) FindBy(fn stream.PredicateFn) *Pipeline {
	return pipeline.Apply(transformers.FindBy(fn))
}

func (pipeline *Pipeline) TakeFirst(n int) *Pipeline {
	return pipeline.Apply(transformers.TakeFirst(n))
}

func (pipeline *Pipeline) Take(fn stream.PredicateFn) *Pipeline {
	return pipeline.Filter(fn)
}

func (pipeline *Pipeline) DropFirst(n int) *Pipeline {
	return pipeline.Apply(transformers.DropFirst(n))
}

func (pipeline *Pipeline) Drop(fn stream.PredicateFn) *Pipeline {
	return pipeline.Take(func(data stream.T) bool { return !fn(data) })
}

func (pipeline *Pipeline) Reduce(acc stream.T, fn stream.ReduceFn) *Pipeline {
	return pipeline.Apply(transformers.Reduce(acc, fn))
}

func (pipeline *Pipeline) Flatten() *Pipeline {
	return pipeline.ApplyParallel(transformers.Flatten())
}

func (pipeline *Pipeline) Batch(size int) *Pipeline {
	return pipeline.Apply(transformers.Batch(size))
}

func (pipeline *Pipeline) BatchBy(batch stream.Batch) *Pipeline {
	return pipeline.Apply(transformers.BatchBy(batch))
}

func (pipeline *Pipeline) Then(consumer stream.Consumer) error {
	consumer.Attach(pipeline.Context)
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

func (pipeline *Pipeline) SortBy(fn stream.SortByFn) ([]stream.T, error) {
	items, err := pipeline.Collect()

	if err != nil {
		return nil, err
	}

	fn.Sort(items)
	return items, nil
}

func (pipeline *Pipeline) GroupBy(groupFn stream.MapFn) (stream.Groups, error) {
	result := make(stream.Groups)
	return result, pipeline.Then(consumers.GroupBy(groupFn, result))
}

func (pipeline *Pipeline) Count() (int, error) {
	items, err := pipeline.Collect()
	return len(items), err
}

func (pipeline *Pipeline) Drain() error {
	return pipeline.Then(consumers.Drainer())
}
