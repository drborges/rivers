package dispatchers

import "github.com/drborges/rivers/rx"

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) If(fn rx.PredicateFn) rx.Dispatcher {
	return &ifDispatcher{
		context: b.context,
		fn:      fn,
	}
}

func (b *Builder) Always() rx.Dispatcher {
	return &ifDispatcher{
		context: b.context,
		fn:      func(_ rx.T) bool { return true },
	}
}
