package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/rx"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type batch struct {
	size  int
	items []rx.T
}

func (batch *batch) Full() bool {
	return len(batch.items) == batch.size
}

func (batch *batch) Empty() bool {
	return len(batch.items) == 0
}

func (batch *batch) Commit(out rx.OutStream) {
	out <- batch.items
	batch.items = []rx.T{}
}

func (batch *batch) Add(data rx.T) {
	batch.items = append(batch.items, data)
}

func TestBatcher(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(3)
			out <- 1
			out <- 2
			out <- 3
			close(out)

			Convey("When I apply the batch transformer to the stream", func() {
				next := transformers.New(context).Batch(2).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{[]rx.T{1, 2}, []rx.T{3}})
				})
			})

			Convey("When I apply the batch by transformer to the stream", func() {
				next := transformers.New(context).BatchBy(&batch{size: 1}).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{[]rx.T{1}, []rx.T{2}, []rx.T{3}})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).Flatten().Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
