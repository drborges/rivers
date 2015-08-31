package producers

import "github.com/drborges/rivers/stream"

type fromRange struct {
	context stream.Context
	from    int
	To      int
}

func (p *fromRange) Produce() stream.Readable {
	capacity := p.To - p.from
	if capacity < 0 {
		capacity = 0
	}

	reader, writer := stream.New(capacity)
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
