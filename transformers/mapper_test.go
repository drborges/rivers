package transformers_test

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMapper(t *testing.T) {
	inc := func(d rx.T) rx.T { return d.(int) + 1 }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply a mapper transformer to the stream", func() {
				transformed := transformers.New(context).Map(inc).Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(transformed.Read(), ShouldResemble, []rx.T{2, 3})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					next := transformers.New(context).Map(inc).Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
