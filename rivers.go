package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	"os"
)

type producer interface {
	From(producer stream.Producer) *Stage
	FromStream(readable stream.Readable) *Stage
	FromRange(from, to int) *Stage
	FromData(data ...stream.T) *Stage
	FromSlice(slice stream.T) *Stage
	FromSocket(protocol, addr string) *Stage
	FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *Stage
	Combine(in ...stream.Readable) *Stage
	CombineZipping(in ...stream.Readable) *Stage
	CombineZippingBy(fn stream.ReduceFn, in ...stream.Readable) *Stage
}

type Stage struct {
	in           stream.Readable
	context      stream.Context
	producers    *producers.Builder
	consumers    *consumers.Builder
	combiners    *combiners.Builder
	dispatchers  *dispatchers.Builder
	transformers *transformers.Builder
}

func (s *Stage) NewFrom(in stream.Readable) *Stage {
	return &Stage{
		in:           in,
		context:      s.context,
		producers:    s.producers,
		consumers:    s.consumers,
		combiners:    s.combiners,
		dispatchers:  s.dispatchers,
		transformers: s.transformers,
	}
}

func New() producer {
	return NewWith(NewContext())
}

func NewWith(context stream.Context) producer {
	return &Stage{
		context:      context,
		producers:    producers.New(context),
		consumers:    consumers.New(context),
		combiners:    combiners.New(context),
		dispatchers:  dispatchers.New(context),
		transformers: transformers.New(context),
	}
}

func (stage *Stage) Split() (*Stage, *Stage) {
	stages := stage.SplitN(2)
	return stages[0], stages[1]
}

func (stage *Stage) SplitN(n int) []*Stage {
	stages := make([]*Stage, n)
	streams := make([]stream.Writable, n)
	for i := 0; i < n; i++ {
		in, out := stream.New(cap(stage.in))
		stages[i] = stage.NewFrom(in)
		streams[i] = out
	}
	stage.dispatchers.Always().Dispatch(stage.in, streams...)
	return stages
}

func (stage *Stage) Partition(fn stream.PredicateFn) (*Stage, *Stage) {
	lhsIn, lhsOut := stream.New(cap(stage.in))
	rhsIn := stage.dispatchers.If(fn).Dispatch(stage.in, lhsOut)

	return stage.NewFrom(lhsIn), stage.NewFrom(rhsIn)
}

func (stage *Stage) Combine(in ...stream.Readable) *Stage {
	return stage.NewFrom(stage.combiners.FIFO().Combine(in...))
}

func (stage *Stage) CombineZipping(in ...stream.Readable) *Stage {
	return stage.NewFrom(stage.combiners.Zip().Combine(in...))
}

func (stage *Stage) CombineZippingBy(fn stream.ReduceFn, in ...stream.Readable) *Stage {
	return stage.NewFrom(stage.combiners.ZipBy(fn).Combine(in...))
}

func (stage *Stage) Dispatch(out ...stream.Writable) *Stage {
	return stage.NewFrom(stage.dispatchers.Always().Dispatch(stage.in, out...))
}

func (stage *Stage) DispatchIf(fn stream.PredicateFn, out ...stream.Writable) *Stage {
	return stage.NewFrom(stage.dispatchers.If(fn).Dispatch(stage.in, out...))
}

func (stage *Stage) From(producer stream.Producer) *Stage {
	return stage.NewFrom(producer.Produce())
}

func (stage *Stage) FromStream(readable stream.Readable) *Stage {
	return stage.NewFrom(readable)
}

func (stage *Stage) FromRange(from, to int) *Stage {
	return stage.NewFrom(stage.producers.FromRange(from, to).Produce())
}

func (stage *Stage) FromFileByLine(file *os.File) *Stage {
	return stage.NewFrom(stage.producers.FromFile(file).ByLine().Produce())
}

func (stage *Stage) FromFileByDelimiter(file *os.File, delimiter byte) *Stage {
	return stage.NewFrom(stage.producers.FromFile(file).ByDelimiter(delimiter).Produce())
}

func (stage *Stage) FromData(data ...stream.T) *Stage {
	return stage.NewFrom(stage.producers.FromData(data...).Produce())
}

func (stage *Stage) FromSlice(slice stream.T) *Stage {
	return stage.NewFrom(stage.producers.FromSlice(slice).Produce())
}

func (stage *Stage) FromSocket(protocol, addr string) *Stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanners.NewLineScanner()).Produce())
}

func (stage *Stage) FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *Stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanner).Produce())
}

func (stage *Stage) Apply(t stream.Transformer) *Stage {
	return stage.NewFrom(t.Transform(stage.in))
}

func (stage *Stage) Filter(fn stream.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.Filter(fn))
}

func (stage *Stage) ProcessWith(fn stream.OnDataFn) *Stage {
	return stage.Apply(stage.transformers.ProcessWith(fn))
}

func (stage *Stage) Map(fn stream.MapFn) *Stage {
	return stage.Apply(stage.transformers.Map(fn))
}

func (stage *Stage) Each(fn stream.EachFn) *Stage {
	return stage.Apply(stage.transformers.Each(fn))
}

func (stage *Stage) FindBy(fn stream.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.FindBy(fn))
}

func (stage *Stage) Take(n int) *Stage {
	return stage.Apply(stage.transformers.Take(n))
}

func (stage *Stage) TakeIf(fn stream.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.TakeIf(fn))
}

func (stage *Stage) DropIf(fn stream.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.DropIf(fn))
}

func (stage *Stage) Reduce(acc stream.T, fn stream.ReduceFn) *Stage {
	return stage.Apply(stage.transformers.Reduce(acc, fn))
}

func (stage *Stage) Flatten() *Stage {
	return stage.Apply(stage.transformers.Flatten())
}

func (stage *Stage) SortBy(fn stream.SortByFn) *Stage {
	return stage.Apply(stage.transformers.SortBy(fn))
}

func (stage *Stage) Batch(size int) *Stage {
	return stage.Apply(stage.transformers.Batch(size))
}

func (stage *Stage) BatchBy(batch stream.Batch) *Stage {
	return stage.Apply(stage.transformers.BatchBy(batch))
}

func (stage *Stage) Sink() stream.Readable {
	return stage.in
}

func (stage *Stage) Collect() ([]stream.T, error) {
	var data []stream.T
	return data, stage.CollectAs(&data)
}

func (stage *Stage) CollectAs(data interface{}) error {
	stage.consumers.ItemsCollector(data).Consume(stage.in)
	return stage.context.Err()
}

func (stage *Stage) CollectFirst() (stream.T, error) {
	var data stream.T
	return data, stage.CollectFirstAs(&data)
}

func (stage *Stage) CollectFirstAs(data interface{}) error {
	stage.consumers.LastItemCollector(data).Consume(stage.Take(1).Sink())
	return stage.context.Err()
}

func (stage *Stage) CollectLast() (stream.T, error) {
	var data stream.T
	return data, stage.CollectLastAs(&data)
}

func (stage *Stage) CollectLastAs(data interface{}) error {
	stage.consumers.LastItemCollector(data).Consume(stage.in)
	return stage.context.Err()
}

func (stage *Stage) Drain() error {
	stage.consumers.Drainer().Consume(stage.in)
	return stage.context.Err()
}
