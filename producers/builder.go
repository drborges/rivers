package producers

import (
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/scanners"
	"os"
)

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) FromRange(from, to int) stream.Producer {
	return &fromRange{
		context: b.context,
		from:    from,
		To:      to,
	}
}

func (b *Builder) FromSlice(slice stream.T) stream.Producer {
	return &fromSlice{
		context: b.context,
		slice:   slice,
	}
}

func (b *Builder) FromData(data ...stream.T) stream.Producer {
	return &fromSlice{
		context: b.context,
		slice:   data,
	}
}

func (b *Builder) FromFile(f *os.File) *fromFile {
	return &fromFile{b.context, f}
}

func (b *Builder) FromSocket(protocol, addr string, scanner scanners.Scanner) stream.Producer {
	return &fromSocket{
		context:  b.context,
		protocol: protocol,
		addr:     addr,
		scanner:  scanner,
	}
}
