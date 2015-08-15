package rivers

import (
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/rx"
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/transformers"
	"os"
)

type producer interface {
	From(producer rx.Producer) *Stage
	FromStream(stream rx.InStream) *Stage
	FromRange(from, to int) *Stage
	FromData(data ...rx.T) *Stage
	FromSlice(slice rx.T) *Stage
	FromSocket(protocol, addr string) *Stage
	FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *Stage
	Combine(in ...rx.InStream) *Stage
	CombineZipping(in ...rx.InStream) *Stage
	CombineZippingBy(fn rx.ReduceFn, in ...rx.InStream) *Stage
}

type Stage struct {
	in           rx.InStream
	context      rx.Context
	producers    *producers.Builder
	consumers    *consumers.Builder
	combiners    *combiners.Builder
	dispatchers  *dispatchers.Builder
	transformers *transformers.Builder
}

func (s *Stage) NewFrom(in rx.InStream) *Stage {
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

func NewWith(context rx.Context) producer {
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
	lhsIn, lhsOut := rx.NewStream(cap(stage.in))
	rhsIn, rhsOut := rx.NewStream(cap(stage.in))
	stage.dispatchers.Always().Dispatch(stage.in, lhsOut, rhsOut)

	return stage.NewFrom(lhsIn), stage.NewFrom(rhsIn)
}

func (stage *Stage) Partition(fn rx.PredicateFn) (*Stage, *Stage) {
	lhsIn, lhsOut := rx.NewStream(cap(stage.in))
	rhsIn := stage.dispatchers.If(fn).Dispatch(stage.in, lhsOut)

	return stage.NewFrom(lhsIn), stage.NewFrom(rhsIn)
}

func (stage *Stage) Combine(in ...rx.InStream) *Stage {
	return stage.NewFrom(stage.combiners.FIFO().Combine(in...))
}

func (stage *Stage) CombineZipping(in ...rx.InStream) *Stage {
	return stage.NewFrom(stage.combiners.Zip().Combine(in...))
}

func (stage *Stage) CombineZippingBy(fn rx.ReduceFn, in ...rx.InStream) *Stage {
	return stage.NewFrom(stage.combiners.ZipBy(fn).Combine(in...))
}

func (stage *Stage) Dispatch(out ...rx.OutStream) *Stage {
	return stage.NewFrom(stage.dispatchers.Always().Dispatch(stage.in, out...))
}

func (stage *Stage) DispatchIf(fn rx.PredicateFn, out ...rx.OutStream) *Stage {
	return stage.NewFrom(stage.dispatchers.If(fn).Dispatch(stage.in, out...))
}

func (stage *Stage) From(producer rx.Producer) *Stage {
	return stage.NewFrom(producer.Produce())
}

func (stage *Stage) FromStream(stream rx.InStream) *Stage {
	return stage.NewFrom(stream)
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

func (stage *Stage) FromData(data ...rx.T) *Stage {
	return stage.NewFrom(stage.producers.FromData(data...).Produce())
}

func (stage *Stage) FromSlice(slice rx.T) *Stage {
	return stage.NewFrom(stage.producers.FromSlice(slice).Produce())
}

func (stage *Stage) FromSocket(protocol, addr string) *Stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanners.NewLineScanner()).Produce())
}

func (stage *Stage) FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *Stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanner).Produce())
}

func (stage *Stage) Apply(t rx.Transformer) *Stage {
	return stage.NewFrom(t.Transform(stage.in))
}

func (stage *Stage) Filter(fn rx.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.Filter(fn))
}

func (stage *Stage) ProcessWith(fn rx.OnDataFn) *Stage {
	return stage.Apply(stage.transformers.ProcessWith(fn))
}

func (stage *Stage) Map(fn rx.MapFn) *Stage {
	return stage.Apply(stage.transformers.Map(fn))
}

func (stage *Stage) Each(fn rx.EachFn) *Stage {
	return stage.Apply(stage.transformers.Each(fn))
}

func (stage *Stage) FindBy(fn rx.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.FindBy(fn))
}

func (stage *Stage) TakeBy(fn rx.PredicateFn) *Stage {
	return stage.Apply(stage.transformers.TakeBy(fn))
}

func (stage *Stage) Reduce(acc rx.T, fn rx.ReduceFn) *Stage {
	return stage.Apply(stage.transformers.Reduce(acc, fn))
}

func (stage *Stage) Flatten() *Stage {
	return stage.Apply(stage.transformers.Flatten())
}

func (stage *Stage) SortBy(fn rx.SortByFn) *Stage {
	return stage.Apply(stage.transformers.SortBy(fn))
}

func (stage *Stage) Batch(size int) *Stage {
	return stage.Apply(stage.transformers.Batch(size))
}

func (stage *Stage) BatchBy(batch rx.Batch) *Stage {
	return stage.Apply(stage.transformers.BatchBy(batch))
}

func (stage *Stage) Sink() rx.InStream {
	return stage.in
}

func (stage *Stage) Drain() error {
	stage.consumers.Drainer().Consume(stage.in)
	return stage.context.Err()
}
