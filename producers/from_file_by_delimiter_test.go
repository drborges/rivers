package producers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/producers"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
	"os"
	"io/ioutil"
)

func TestFromFileByDelimiter(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a file with some data", func() {
			ioutil.WriteFile("/tmp/from_file_by_line", []byte("Hello there folks!"), 0644)
			file, _ := os.Open("/tmp/from_file_by_line")

			Convey("When I produce data from the file", func() {
				stream := producers.New(context).FromFile(file).ByDelimiter(' ').Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(stream.Read(), ShouldResemble, []rx.T{"Hello", "there", "folks!"})
				})
			})
		})
	})
}