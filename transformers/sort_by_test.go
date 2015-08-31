package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSortBy(t *testing.T) {
	sorter := func(a, b stream.T) bool { return a.(int) < b.(int) }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(3)
			out <- 3
			out <- 2
			out <- 1
			close(out)

			Convey("When I apply a mapper transformer to the stream", func() {
				next := transformers.New(context).SortBy(sorter).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []stream.T{1, 2, 3})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).SortBy(sorter).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
