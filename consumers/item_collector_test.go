package consumers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/rx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestItemCollector(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the collector consumer", func() {
				var number int
				consumers.New(context).ItemCollector(&number).Consume(in)

				Convey("Then data is collected out of the stream", func() {
					So(number, ShouldResemble, 1)

					data, opened := <-in
					So(data, ShouldEqual, 2)
					So(opened, ShouldBeTrue)
				})
			})

			Convey("When I apply the collector consuming data into a non pointer", func() {
				var number int
				collect := func() {
					consumers.New(context).ItemCollector(number)
				}

				Convey("Then it panics", func() {
					So(collect, ShouldPanicWith, rx.ErrNoSuchPointer)
				})
			})
		})
	})
}
