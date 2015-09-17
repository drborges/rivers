package dispatchers

import "github.com/drborges/rivers/stream"

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) If(fn stream.PredicateFn) stream.Dispatcher {
	return &dispatcher{
		context: b.context,
		fn:      fn,
	}
}

func (b *Builder) Always() stream.Dispatcher {
	return &dispatcher{
		context: b.context,
		fn:      func(_ stream.T) bool { return true },
	}
}
