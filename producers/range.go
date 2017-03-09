package producers

import (
	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/pipeline"
	"github.com/drborges/rivers/stream"
)

// Range creates a producer that generates integers within the given range.
func Range(from, to int) pipeline.Producer {
	return func(ctx ctxtree.Context) stream.Reader {
		r, w := stream.NewWithContext(ctx)

		go func() {
			defer w.Close(nil)

			for i := from; i <= to; i++ {
				if err := w.Write(i); err != nil {
					return
				}
			}
		}()

		return r
	}
}
