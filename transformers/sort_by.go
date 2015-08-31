package transformers

import "github.com/drborges/rivers/rx"

type sortBy struct {
	context rx.Context
	sorter  rx.SortByFn
}

// TODO Sort on demand rather than waiting to receive all items
func (t *sortBy) Transform(in rx.Readable) rx.Readable {
	reader, writer := rx.NewStream(cap(in))

	go func() {
		defer t.context.Recover()
		defer close(writer)

		items := []rx.T{}
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
