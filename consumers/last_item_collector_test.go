package consumers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLastItemCollector(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the collector consumer", func() {
				var number int
				consumers.New(context).LastItemCollector(&number).Consume(in)

				Convey("Then data is collected out of the stream", func() {
					So(number, ShouldResemble, 2)

					_, opened := <-in
					So(opened, ShouldBeFalse)
				})
			})

			Convey("When I apply the collector consuming data into a non pointer", func() {
				var number int
				collect := func() {
					consumers.New(context).LastItemCollector(number)
				}

				Convey("Then it panics", func() {
					So(collect, ShouldPanicWith, consumers.ErrNoSuchPointer)
				})
			})
		})
	})
}
