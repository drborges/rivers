package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEach(t *testing.T) {
	collect := func(items *[]stream.T) stream.EachFn {
		return func(data stream.T) {
			*items = append(*items, data)
		}
	}

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				var items []stream.T
				next := transformers.New(context).Each(collect(&items)).Transform(in)

				Convey("Then all items are sent to the next stage", func() {
					So(next.Read(), ShouldResemble, []stream.T{1, 2})

					Convey("And all items are transformed", func() {
						So(items, ShouldResemble, []stream.T{1, 2})
					})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					var items []stream.T
					next := transformers.New(context).Each(collect(&items)).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)

						Convey("And no item is transformed", func() {
							So(items, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
