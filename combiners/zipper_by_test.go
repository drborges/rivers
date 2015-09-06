package combiners_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestZipperBy(t *testing.T) {
	adder := func(a, b stream.T) stream.T {
		return a.(int) + b.(int)
	}

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in1, out1 := stream.New(2)
			out1 <- 1
			out1 <- 2
			close(out1)

			in2, out2 := stream.New(4)
			out2 <- 3
			out2 <- 4
			out2 <- 5
			out2 <- 6
			close(out2)

			Convey("When I apply the combiner to the streams", func() {
				combiner := combiners.ZipBy(adder)
				combiner.(stream.Bindable).Bind(context)
				combined := combiner.Combine(in1, in2)

				Convey("Then a transformed stream is returned", func() {
					So(combined.ReadAll(), ShouldResemble, []stream.T{4, 6, 5, 6})
				})
			})

			Convey("When I close the context", func() {
				context.Close(nil)

				Convey("And I apply the transformer to the stream", func() {
					combiner := combiners.ZipBy(adder)
					combiner.(stream.Bindable).Bind(context)
					combined := combiner.Combine(in1, in2)

					Convey("Then no item is sent to the next stage", func() {
						So(combined.ReadAll(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
