package producers

import (
	"github.com/drborges/riversv2/rx"
	"github.com/drborges/riversv2/scanners"
	"net"
)

type fromSocket struct {
	context  rx.Context
	protocol string
	addr     string
	scanner  scanners.Scanner
}

func (socket *fromSocket) Produce() rx.InStream {
	reader, writer := rx.NewStream(100)

	go func() {
		defer socket.context.Recover()
		defer close(writer)

		conn, err := net.Dial(socket.protocol, socket.addr)
		if err != nil {
			return
		}

		socket.scanner.Attach(conn)
		for {
			select {
			case <-socket.context.Closed():
				return
			default:
				if message, err := socket.scanner.Scan(); err == nil {
					writer <- message
					continue
				}
				return
			}
		}
	}()

	return reader
}
