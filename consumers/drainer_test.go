package consumers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/consumers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDrainer(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the drainer consumer", func() {
				consumer := consumers.Drainer()
				consumer.Attach(context)
				consumer.Consume(in)

				Convey("Then the stream is drained", func() {
					data, opened := <-in
					So(data, ShouldBeNil)
					So(opened, ShouldBeFalse)
				})
			})
		})
	})
}
