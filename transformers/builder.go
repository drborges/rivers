package transformers

import "github.com/drborges/riversv2/rx"

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Filter(fn rx.PredicateFn) rx.Transformer {
	return &filter{
		context: b.context,
		predicate: fn,
	}
}

func (b *Builder) FindBy(fn rx.PredicateFn) rx.Transformer {
	return &findBy{
		context: b.context,
		predicate: fn,
	}
}

func (b *Builder) Map(fn rx.MapFn) rx.Transformer {
	return &mapper{
		context: b.context,
		mapFn: fn,
	}
}

func (b *Builder) Reduce(acc rx.T, fn rx.ReduceFn) rx.Transformer {
	return &reducer{
		context: b.context,
		reduceFn: fn,
		acc: acc,
	}
}

func (b *Builder) Flatten() rx.Transformer {
	return &flatten{
		context: b.context,
	}
}

func (b *Builder) Batch(size int) rx.Transformer {
	return &batcher{
		context: b.context,
		size: size,
	}
}

func (b *Builder) SortBy(sorter rx.SortByFn) rx.Transformer {
	return &sortBy{
		context: b.context,
		sorter: sorter,
	}
}

func (b *Builder) Each(fn rx.HandleFn) rx.Transformer {
	return &each{
		context: b.context,
		handler: fn,
	}
}
