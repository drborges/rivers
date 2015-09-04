package transformers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/transformers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFindBy(t *testing.T) {
	evens := func(d stream.T) bool { return d.(int)%2 == 0 }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(3)
			out <- 1
			out <- 2
			out <- 4
			close(out)

			Convey("When I apply the transformer to the stream", func() {
				transformer := transformers.FindBy(evens)
				transformer.(stream.Bindable).Bind(context)
				next := transformer.Transform(in)

				Convey("Then a transformed stream is returned", func() {
					So(next.Read(), ShouldResemble, []stream.T{2})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					transformer := transformers.FindBy(evens)
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
