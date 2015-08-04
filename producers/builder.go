package producers

import (
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/scanners"
	"os"
)

type Builder struct {
	context rx.Context
}

func New(c rx.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) FromRange(from, to int) rx.Producer {
	return &fromRange{
		context: b.context,
		from:    from,
		To:      to,
	}
}

func (b *Builder) FromSlice(slice rx.T) rx.Producer {
	return &fromSlice{
		context: b.context,
		slice:   slice,
	}
}

func (b *Builder) FromData(data ...rx.T) rx.Producer {
	return &fromSlice{
		context: b.context,
		slice:   data,
	}
}

func (b *Builder) FromFile(f *os.File) *fromFile {
	return &fromFile{b.context, f}
}

func (b *Builder) FromSocket(protocol, addr string, scanner scanners.Scanner) rx.Producer {
	return &fromSocket{
		context:  b.context,
		protocol: protocol,
		addr:     addr,
		scanner:  scanner,
	}
}
