package consumers

import "github.com/drborges/riversv2/rx"

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Drainer() rx.Consumer {
	return &drainer{b.context}
}
