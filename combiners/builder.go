package combiners

import "github.com/drborges/riversv2/rx"

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Zip() rx.Combiner {
	return &zip{
		context: b.context,
	}
}

func (b *Builder) ZipBy(fn rx.ReduceFn) rx.Combiner {
	return &zipBy{
		context: b.context,
		fn:      fn,
	}
}

func (b *Builder) FIFO() rx.Combiner {
	return &fifo{
		context: b.context,
	}
}
