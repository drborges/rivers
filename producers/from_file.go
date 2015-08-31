package producers

import (
	"github.com/drborges/rivers/stream"
	"os"
)

type fromFile struct {
	context stream.Context
	file    *os.File
}

func (builder *fromFile) ByLine() stream.Producer {
	return &fromFileByLine{builder.context, builder.file}
}

func (builder *fromFile) ByDelimiter(delimiter byte) stream.Producer {
	return &fromFileByDelimiter{builder.context, builder.file, delimiter}
}
