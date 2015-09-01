package transformers

import "github.com/drborges/rivers/stream"

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) Filter(fn stream.PredicateFn) stream.Transformer {
	return &filter{
		context:   b.context,
		predicate: fn,
	}
}

func (b *Builder) FindBy(fn stream.PredicateFn) stream.Transformer {
	return &findBy{
		context:   b.context,
		predicate: fn,
	}
}

func (b *Builder) Take(n int) stream.Transformer {
	return &takeN{
		context: b.context,
		n:       n,
	}
}

func (b *Builder) TakeIf(fn stream.PredicateFn) stream.Transformer {
	return b.Filter(fn)
}

func (b *Builder) DropIf(fn stream.PredicateFn) stream.Transformer {
	return b.Filter(func(data stream.T) bool { return !fn(data) })
}

func (b *Builder) Map(fn stream.MapFn) stream.Transformer {
	return &mapper{
		context: b.context,
		mapFn:   fn,
	}
}

func (b *Builder) OnData(fn stream.OnDataFn) stream.Transformer {
	return &processor{
		context: b.context,
		onData:  fn,
	}
}

func (b *Builder) Reduce(acc stream.T, fn stream.ReduceFn) stream.Transformer {
	return &reducer{
		context:    b.context,
		reduceFn:   fn,
		initialAcc: acc,
	}
}

func (b *Builder) Flatten() stream.Transformer {
	return &flatten{
		context: b.context,
	}
}

func (b *Builder) Batch(size int) stream.Transformer {
	return &batcher{
		context: b.context,
		batch:   &batch{size: size},
	}
}

func (b *Builder) BatchBy(batch stream.Batch) stream.Transformer {
	return &batcher{
		context: b.context,
		batch:   batch,
	}
}

func (b *Builder) SortBy(sorter stream.SortByFn) stream.Transformer {
	return &sortBy{
		context: b.context,
		sorter:  sorter,
	}
}

func (b *Builder) Each(fn stream.EachFn) stream.Transformer {
	return &each{
		context: b.context,
		handler: fn,
	}
}
