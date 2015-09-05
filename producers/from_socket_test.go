package producers_test

import (
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/producers"
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/stream"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestFromSocket(t *testing.T) {
	listen := func() (net.Listener, string) {
		port := ":8484"
		ln, err := net.Listen("tcp", port)
		if err != nil {
			port = ":8585"
			ln, _ = net.Listen("tcp", port)
		}
		return ln, port
	}

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a tcp server accepting connections to send data", func() {
			ln, addr := listen()

			go func() {
				conn, _ := ln.Accept()
				defer conn.Close()
				conn.Write([]byte("Hello there\n"))
				conn.Write([]byte("rivers!\n"))
				conn.Write([]byte("super cool!\n"))
			}()

			Convey("When I produce a data stream from a tcp connection scanning by line", func() {
				producer := producers.FromSocket("tcp", addr, scanners.NewLineScanner())
				producer.(stream.Bindable).Bind(context)
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.ReadAll(), ShouldResemble, []stream.T{
						[]byte("Hello there"),
						[]byte("rivers!"),
						[]byte("super cool!"),
					})
				})
			})

			Convey("When I produce a data stream from a tcp connection scanning by word", func() {
				producer := producers.FromSocket("tcp", addr, scanners.NewWordScanner())
				producer.(stream.Bindable).Bind(context)
				readable := producer.Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(readable.ReadAll(), ShouldResemble, []stream.T{
						[]byte("Hello"),
						[]byte("there"),
						[]byte("rivers!"),
						[]byte("super"),
						[]byte("cool!"),
					})
				})
			})
		})
	})
}
