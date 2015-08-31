package consumers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestItemsCollector(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the collector consumer", func() {
				var data []stream.T
				consumers.New(context).ItemsCollector(&data).Consume(in)

				Convey("Then data is collected out of the stream", func() {
					So(data, ShouldResemble, []stream.T{1, 2})

					data, opened := <-in
					So(data, ShouldBeNil)
					So(opened, ShouldBeFalse)
				})
			})

			Convey("When I apply the collector consuming data into a non slice pointer", func() {
				var data []stream.T
				collect := func() {
					consumers.New(context).ItemsCollector(data)
				}

				Convey("Then it panics", func() {
					So(collect, ShouldPanicWith, consumers.ErrNoSuchSlicePointer)
				})
			})
		})
	})
}
