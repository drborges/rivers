package transformers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
)

func TestFlatten(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(3)
			out <- []rx.T{1, 2}
			out <- []rx.T{3}
			out <- 4
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				next := transformers.New(context).Flatten().Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{1, 2, 3, 4})
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
