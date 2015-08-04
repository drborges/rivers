package consumers_test

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/consumers"
	"github.com/drborges/riversv2/rx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDrainer(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := rx.NewStream(2)
			out <- 1
			out <- 2
			close(out)

			Convey("When I apply the drainer consumer", func() {
				consumers.New(context).Drainer().Consume(in)

				Convey("Then the stream is drained", func() {
					data, opened := <-in
					So(data, ShouldBeNil)
					So(opened, ShouldBeFalse)
				})
			})
		})
	})
}
