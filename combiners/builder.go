package combiners

import "github.com/drborges/rivers/stream"

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Zip() stream.Combiner {
	return &zip{
		context: b.context,
	}
}

func (b *Builder) ZipBy(fn stream.ReduceFn) stream.Combiner {
	return &zipBy{
		context: b.context,
		fn:      fn,
	}
}

func (b *Builder) FIFO() stream.Combiner {
	return &fifo{
		context: b.context,
	}
}
