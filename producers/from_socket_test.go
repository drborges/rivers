package producers_test

import (
	"testing"
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/producers"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/drborges/riversv2/rx"
	"net"
	"github.com/drborges/riversv2/scanners"
)

func TestFromSocket(t *testing.T) {
	listen := func() (net.Listener, string) {
		port := ":8080"
		ln, err := net.Listen("tcp", port)
		if err != nil {
			port = ":8081"
			ln, _ = net.Listen("tcp", port)
		}
		return ln, port
	}

	Convey("Given I have a context", t, func() {
		context := rivers.NewContext()

		Convey("And I have a tcp server accepting connections to send data", func() {
			ln, port := listen()

			go func() {
				conn, _ := ln.Accept()
				defer conn.Close()
				conn.Write([]byte("Hello there\n"))
				conn.Write([]byte("rivers!\n"))
				conn.Write([]byte("super cool!\n"))
			}()

			Convey("When I produce a data stream from a tcp connection scanning by line", func() {
				stream := producers.New(context).FromSocket("tcp", port, scanners.NewLineScanner()).Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(stream.Read(), ShouldResemble, []rx.T{
						[]byte("Hello there"),
						[]byte("rivers!"),
						[]byte("super cool!"),
					})
				})
			})

			Convey("When I produce a data stream from a tcp connection scanning by word", func() {
				stream := producers.New(context).FromSocket("tcp", port, scanners.NewWordScanner()).Produce()

				Convey("Then I can read the produced data from the stream", func() {
					So(stream.Read(), ShouldResemble, []rx.T{
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