package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFlatten(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(3)
			out <- []stream.T{1, 2}
			out <- []stream.T{3}
			out <- 4
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				transformer := transformers.Flatten()
				transformer.Attach(context)
				next := transformer.Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.ReadAll(), ShouldResemble, []stream.T{1, 2, 3, 4})
				})
			})

			Convey("When I close the context", func() {
				context.Close(stream.Done)

				Convey("And I apply the transformer to the stream", func() {
					transformer := transformers.Flatten()
					transformer.Attach(context)
					next := transformer.Transform(in)

					Convey("Then no item is sent to the next stage", func() {
						So(next.ReadAll(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
