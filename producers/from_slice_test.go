package producers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/producers"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
)

func TestFromSlice(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a slice producer", func() {
			numbers := []int{1, 2, 3}
			producer := producers.New(context).FromSlice(numbers)

			Convey("When I produce data", func() {
				stream := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(stream.Read(), ShouldResemble, []rx.T{1, 2, 3})
				})
			})
		})

		Convey("And I have a data producer", func() {
			producer := producers.New(context).FromData(1, 2, 3)

			Convey("When I produce data", func() {
				stream := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(stream.Read(), ShouldResemble, []rx.T{1, 2, 3})
				})
			})
		})
	})
}