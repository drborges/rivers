package transformers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
)

func TestReducer(t *testing.T) {
	sum := func(acc, next rx.T) rx.T { return acc.(int) + next.(int) }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(3)
			out <- 1
			out <- 2
			out <- 3
			close(out)

			Convey("When I apply a mapper transformer to the stream", func() {
				next := transformers.New(context).Reduce(0, sum).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{6})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).Reduce(0, sum).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
