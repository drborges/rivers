package producers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFromSlice(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a slice producer", func() {
			numbers := []int{1, 2, 3}
			producer := producers.FromSlice(numbers)

			Convey("When I produce data", func() {
				producer.(stream.Bindable).Bind(context)
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.Read(), ShouldResemble, []stream.T{1, 2, 3})
				})
			})
		})

		Convey("And I have a data producer", func() {
			producer := producers.FromData(1, 2, 3)

			Convey("When I produce data", func() {
				producer.(stream.Bindable).Bind(context)
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.Read(), ShouldResemble, []stream.T{1, 2, 3})
				})
			})
		})
	})
}
