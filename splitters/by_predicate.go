package splitters

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
)

func ByPredicate(fn rivers.Predicate) rivers.Splitter {
	return func(upstream stream.Reader) (stream.Reader, stream.Reader) {
		reader1, writer1 := upstream.NewDownstream()
		reader2, writer2 := upstream.NewDownstream()

		go func() {
			defer writer1.Close(nil)
			defer writer2.Close(nil)

			for data := range upstream.Read() {
				if fn(data) {
					if err := writer1.Write(data); err != nil {
						continue
					}
				} else {
					if err := writer2.Write(data); err != nil {
						continue
					}
				}
			}
		}()

		return reader1, reader2
	}
}
