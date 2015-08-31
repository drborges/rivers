package producers

import "github.com/drborges/rivers/rx"

type fromRange struct {
	context rx.Context
	from    int
	To      int
}

func (p *fromRange) Produce() rx.Readable {
	capacity := p.To - p.from
	if capacity < 0 {
		capacity = 0
	}

	reader, writer := rx.NewStream(capacity)
	go func() {
		defer p.context.Recover()
		defer close(writer)

		for i := p.from; i <= p.To; i++ {
			select {
			case <-p.context.Closed():
				return
			default:
				writer <- i
			}
		}
	}()

	return reader
}
