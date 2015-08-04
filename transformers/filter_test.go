package transformers_test

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFilter(t *testing.T) {
	evens := func(d rx.T) bool { return d.(int)%2 == 0 }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				transformed := transformers.New(context).Filter(evens).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(transformed.Read(), ShouldResemble, []rx.T{2})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).Filter(evens).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
