package transformers

import "github.com/drborges/rivers/stream"

type sortBy struct {
	context stream.Context
	sorter  stream.SortByFn
}

// TODO Sort on demand rather than waiting to receive all items
func (t *sortBy) Transform(in stream.Readable) stream.Readable {
	reader, writer := stream.New(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		items := []stream.T{}
		for {
			select {
			case <-t.context.Closed():
				return
			default:
				data, more := <-in
				if !more {
					t.sorter.Sort(items)
					for _, item := range items {
						writer <- item
					}
					return
				}
				items = append(items, data)
			}
		}
	}()

	return reader
}
