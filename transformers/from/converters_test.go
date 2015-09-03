package from_test

import (
	"encoding/json"
	"github.com/drborges/rivers/transformers/from"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStructToJSON(t *testing.T) {
	type Account struct{ Name string }

	convey.Convey("Given I have an instance of a particular struct", t, func() {
		item := Account{"Diego"}

		convey.Convey("When I convert it to JSON", func() {
			converted := from.StructToJSON(item)

			bytes, _ := json.Marshal(item)
			convey.So(converted, should.Resemble, bytes)
		})
	})
}

func TestJSONToStruct(t *testing.T) {
	type Account struct{ Name string }

	convey.Convey("Given I have an instance of a particular struct", t, func() {
		json := []byte(`{"Name":"Diego"}`)

		convey.Convey("When I convert it to Struct", func() {
			converted := from.JSONToStruct(Account{})(json)
			convey.So(converted, should.Resemble, Account{"Diego"})

			converted = from.JSONToStruct(&Account{})(json)
			convey.So(converted, should.Resemble, &Account{"Diego"})
		})
	})
}
