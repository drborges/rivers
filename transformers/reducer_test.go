package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestReducer(t *testing.T) {
	sum := func(acc, next stream.T) stream.T { return acc.(int) + next.(int) }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(3)
			out <- 1
			out <- 2
			out <- 3
			close(out)

			Convey("When I apply a mapper transformer to the stream", func() {
				transformer := transformers.Reduce(0, sum)
				transformer.(stream.Bindable).Bind(context)
				next := transformer.Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []stream.T{6})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					transformer := transformers.Reduce(0, sum)
					transformer.(stream.Bindable).Bind(context)
					next := transformer.Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
