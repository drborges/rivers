package producers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestFromFileByLine(t *testing.T) {
	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a file with some data", func() {
			ioutil.WriteFile("/tmp/from_file_by_line", []byte("Hello\nthere\nfolks!"), 0644)
			file, _ := os.Open("/tmp/from_file_by_line")

			Convey("When I produce data from the file", func() {
				producer := producers.FromFile(file).ByLine()
				producer.Attach(context)
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.ReadAll(), ShouldResemble, []stream.T{"Hello", "there", "folks!"})
				})
			})
		})
	})
}
