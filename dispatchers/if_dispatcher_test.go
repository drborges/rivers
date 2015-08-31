package dispatchers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/dispatchers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIfDispatcher(t *testing.T) {
	evens := func(d stream.T) bool { return d.(int)%2 == 0 }

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And a stream of data", func() {
			in, out := stream.New(3)
			out <- 2
			out <- 3
			out <- 4
			close(out)

			Convey("When I apply an if dispatcher", func() {
				evensIn, evensOut := stream.New(3)
				sink := dispatchers.New(context).If(evens).Dispatch(in, evensOut)

				Convey("Then items matching the condition are dispatched to the corresponding stream", func() {
					So(evensIn.Read(), ShouldResemble, []stream.T{2, 4})

					Convey("And items not matching the condition are dispatched to the sink stream", func() {
						So(sink.Read(), ShouldResemble, []stream.T{3})
					})
				})
			})

			Convey("When I apply an always dispatcher", func() {
				streamIn1, streamOut1 := stream.New(3)
				streamIn2, streamOut2 := stream.New(2)
				sink := dispatchers.New(context).Always().Dispatch(in, streamOut1, streamOut2)

				Convey("Then all items are dispatched to the corresponding streams", func() {
					So(streamIn1.Read(), ShouldResemble, []stream.T{2, 3, 4})
					So(streamIn2.Read(), ShouldResemble, []stream.T{2, 3, 4})

					Convey("And no item is dispatched to the sink stream", func() {
						So(sink.Read(), ShouldBeEmpty)
					})
				})
			})

			Convey("When I close the context", func() {
				context.Close()

				Convey("And I apply the transformer to the stream", func() {
					evensIn, evensOut := stream.New(2)
					sink := dispatchers.New(context).If(evens).Dispatch(in, evensOut)

					Convey("Then no item is sent to the next stage", func() {
						So(evensIn.Read(), ShouldBeEmpty)
						So(sink.Read(), ShouldBeEmpty)
					})
				})
			})
		})
	})
}
