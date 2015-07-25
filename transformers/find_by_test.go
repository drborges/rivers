package transformers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
)

func TestFindBy(t *testing.T) {
	evens := func(d rx.T) bool { return d.(int) % 2 == 0 }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(3)
			out <- 1
			out <- 2
			out <- 4
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				next := transformers.New(context).FindBy(evens).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []rx.T{2})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).FindBy(evens).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
