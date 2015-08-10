package consumers_test

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/consumers"
	"github.com/drborges/riversv2/rx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDataCollector(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the data collector consumer", func() {
				var data []rx.T
				consumers.New(context).DataCollector(&data).Consume(in)

				Convey("Then data is collected out of the stream", func() {
					So(data, ShouldResemble, []rx.T{1, 2})

					data, opened := <-in
					So(data, ShouldBeNil)
					So(opened, ShouldBeFalse)
				})
			})
		})
	})
}
