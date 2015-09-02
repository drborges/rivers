package where_test

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/drborges/rivers/transformers/where"
	"github.com/smartystreets/assertions/should"
)

func TestItemResemble(t *testing.T) {
	type Account struct{ Name, Email string }

	convey.Convey("Given I have an instance of a particular struct", t, func() {
		item := Account{"Diego", "drborges.cic@gmail.com"}

		convey.Convey("Then it resembles a given subject", func() {
			resemble := where.ItemResemble(Account{item.Name, item.Email})(item)
			convey.So(resemble, should.BeTrue)
		})

		convey.Convey("Then it does not resemble a given subject", func() {
			resemble := where.ItemResemble(Account{"Borges", item.Email})(item)
			convey.So(resemble, should.BeFalse)
		})
	})
}
