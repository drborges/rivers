package combiners_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/combiners"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/smartystreets/assertions/should"
)

func TestFifo(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in1, out1 := stream.New(2)
			out1 <- 1
			out1 <- 2
			close(out1)

			in2, out2 := stream.New(2)
			out2 <- 3
			out2 <- 4
			close(out2)

			Convey("When I apply the combiner to the streams", func() {
				combiner := combiners.FIFO()
				combiner.Attach(context)
				combined := combiner.Combine(in1, in2)

				Convey("Then a transformed stream is returned", func() {
					items := combined.ReadAll()
					So(items, should.Contain, 1)
					So(items, should.Contain, 2)
					So(items, should.Contain, 3)
					So(items, should.Contain, 4)
				})
			})

			Convey("When I close the context", func() {
				context.Close(stream.Done)

				Convey("And I apply the transformer to the stream", func() {
					combiner := combiners.Zip()
					combiner.Attach(context)
					combined := combiner.Combine(in1, in2)

					Convey("Then no item is sent to the next stage", func() {
						So(combined.ReadAll(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
