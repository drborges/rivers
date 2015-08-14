package producers

import (
	"github.com/drborges/rivers/rx"
	"os"
)

type fromFile struct {
	context rx.Context
	file    *os.File
}

func (builder *fromFile) ByLine() rx.Producer {
	return &fromFileByLine{builder.context, builder.file}
}

func (builder *fromFile) ByDelimiter(delimiter byte) rx.Producer {
	return &fromFileByDelimiter{builder.context, builder.file, delimiter}
}
