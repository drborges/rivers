package transformers_test

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBatcher(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(3)
			out <- 1
			out <- 2
			out <- 3
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				next := transformers.New(context).Batch(2).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{[]rx.T{1, 2}, []rx.T{3}})
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
