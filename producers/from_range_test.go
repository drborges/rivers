package producers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/rx"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFromRange(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a range producer", func() {
			producer := producers.New(context).FromRange(1, 3)

			Convey("When I produce data", func() {
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.Read(), ShouldResemble, []rx.T{1, 2, 3})
				})
			})
		})
	})
}
