package where_test

import (
	"github.com/drborges/rivers/transformers/where"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
	"testing"
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

func TestItemIs(t *testing.T) {
	type Account struct{ Name, Email string }

	convey.Convey("Given I have an instance of a particular struct", t, func() {
		item := &Account{"Diego", "drborges.cic@gmail.com"}

		convey.Convey("Then subject is the same instance", func() {
			convey.So(where.ItemIs(item)(item), should.BeTrue)
		})

		convey.Convey("Then subject is not the same instace", func() {
			convey.So(where.ItemIs(&Account{})(item), should.BeFalse)
		})
	})
}

func TestStructHas(t *testing.T) {
	type Account struct{ Name, Email string }

	convey.Convey("Given I have an instance of a particular struct", t, func() {
		item := &Account{"Diego", "drborges.cic@gmail.com"}

		convey.Convey("Then subject matches filter", func() {
			convey.So(where.StructHas("Name", "Diego")(item), should.BeTrue)
		})

		convey.Convey("Then subject is not the same instace", func() {
			convey.So(where.StructHas("Name", "Borges")(item), should.BeFalse)
		})
	})
}
