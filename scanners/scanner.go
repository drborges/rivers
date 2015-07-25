package scanners

import "io"

type Scanner interface {
	Attach(io.Reader)
	Scan() (msg []byte, err error)
}