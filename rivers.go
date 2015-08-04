package rivers

import (
	"github.com/drborges/riversv2/combiners"
	"github.com/drborges/riversv2/consumers"
	"github.com/drborges/riversv2/dispatchers"
	"github.com/drborges/riversv2/producers"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/scanners"
	"github.com/drborges/riversv2/transformers"
	"os"
)

type producer interface {
	From(producer rx.Producer) *stage
	FromStream(stream rx.InStream) *stage
	FromRange(from, to int) *stage
	FromData(data ...rx.T) *stage
	FromSlice(slice rx.T) *stage
	FromSocket(protocol, addr string) *stage
	FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *stage
	Combine(in ...rx.InStream) *stage
	CombineZipping(in ...rx.InStream) *stage
	CombineZippingBy(fn rx.ReduceFn, in ...rx.InStream) *stage
}

type stage struct {
	in           rx.InStream
	producers    *producers.Builder
	consumers    *consumers.Builder
	combiners    *combiners.Builder
	dispatchers  *dispatchers.Builder
	transformers *transformers.Builder
}

func (s *stage) NewFrom(in rx.InStream) *stage {
	return &stage{
		in:           in,
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
	return &stage{
		producers:    producers.New(context),
		consumers:    consumers.New(context),
		combiners:    combiners.New(context),
		dispatchers:  dispatchers.New(context),
		transformers: transformers.New(context),
	}
}

func (stage *stage) Combine(in ...rx.InStream) *stage {
	return stage.NewFrom(stage.combiners.FIFO().Combine(in...))
}

func (stage *stage) CombineZipping(in ...rx.InStream) *stage {
	return stage.NewFrom(stage.combiners.Zip().Combine(in...))
}

func (stage *stage) CombineZippingBy(fn rx.ReduceFn, in ...rx.InStream) *stage {
	return stage.NewFrom(stage.combiners.ZipBy(fn).Combine(in...))
}

func (stage *stage) Dispatch(out ...rx.OutStream) *stage {
	return stage.NewFrom(stage.dispatchers.Always().Dispatch(stage.in, out...))
}

func (stage *stage) DispatchIf(fn rx.PredicateFn, out ...rx.OutStream) *stage {
	return stage.NewFrom(stage.dispatchers.If(fn).Dispatch(stage.in, out...))
}

func (stage *stage) From(producer rx.Producer) *stage {
	return stage.NewFrom(producer.Produce())
}

func (stage *stage) FromStream(stream rx.InStream) *stage {
	return stage.NewFrom(stream)
}

func (stage *stage) FromRange(from, to int) *stage {
	return stage.NewFrom(stage.producers.FromRange(from, to).Produce())
}

func (stage *stage) FromFileByLine(file *os.File) *stage {
	return stage.NewFrom(stage.producers.FromFile(file).ByLine().Produce())
}

func (stage *stage) FromFileByDelimiter(file *os.File, delimiter byte) *stage {
	return stage.NewFrom(stage.producers.FromFile(file).ByDelimiter(delimiter).Produce())
}

func (stage *stage) FromData(data ...rx.T) *stage {
	return stage.NewFrom(stage.producers.FromData(data...).Produce())
}

func (stage *stage) FromSlice(slice rx.T) *stage {
	return stage.NewFrom(stage.producers.FromSlice(slice).Produce())
}

func (stage *stage) FromSocket(protocol, addr string) *stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanners.NewLineScanner()).Produce())
}

func (stage *stage) FromSocketWithScanner(protocol, addr string, scanner scanners.Scanner) *stage {
	return stage.NewFrom(stage.producers.FromSocket(protocol, addr, scanner).Produce())
}

func (stage *stage) Apply(t rx.Transformer) *stage {
	return stage.NewFrom(t.Transform(stage.in))
}

func (stage *stage) Filter(fn rx.PredicateFn) *stage {
	return stage.Apply(stage.transformers.Filter(fn))
}

func (stage *stage) Map(fn rx.MapFn) *stage {
	return stage.Apply(stage.transformers.Map(fn))
}

func (stage *stage) Each(fn rx.HandleFn) *stage {
	return stage.Apply(stage.transformers.Each(fn))
}

func (stage *stage) FindBy(fn rx.PredicateFn) *stage {
	return stage.Apply(stage.transformers.FindBy(fn))
}

func (stage *stage) Reduce(acc rx.T, fn rx.ReduceFn) *stage {
	return stage.Apply(stage.transformers.Reduce(acc, fn))
}

func (stage *stage) Flatten() *stage {
	return stage.Apply(stage.transformers.Flatten())
}

func (stage *stage) SortBy(fn rx.SortByFn) *stage {
	return stage.Apply(stage.transformers.SortBy(fn))
}

func (stage *stage) Batch(size int) *stage {
	return stage.Apply(stage.transformers.Batch(size))
}

func (stage *stage) Sink() rx.InStream {
	return stage.in
}

func (stage *stage) Drain() {
	stage.consumers.Drainer().Consume(stage.in)
}
