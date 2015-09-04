package producers

import (
	"github.com/drborges/rivers/scanners"
	"github.com/drborges/rivers/stream"
	"net"
	"os"
	"reflect"
)

func FromRange(from, to int) stream.Producer {
	return &Observable{
		Capacity: to - from,
		Emit: func(emitter stream.Emitter) {
			for i := from; i <= to; i++ {
				emitter.Emit(i)
			}
		},
	}
}

func FromSlice(slice stream.T) stream.Producer {
	sv := reflect.ValueOf(slice)

	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Ptr {
		panic("No such slice")
	}

	if sv.Kind() == reflect.Ptr && sv.Elem().Kind() != reflect.Slice {
		panic("No such slice")
	}

	return &Observable{
		Capacity: sv.Len(),
		Emit: func(emitter stream.Emitter) {
			for i := 0; i < sv.Len(); i++ {
				emitter.Emit(sv.Index(i).Interface())
			}
		},
	}
}

func FromData(data ...stream.T) stream.Producer {
	return FromSlice(data)
}

func FromFile(f *os.File) *fromFile {
	return &fromFile{f}
}

func FromSocket(protocol, addr string, scanner scanners.Scanner) stream.Producer {
	return &Observable{
		Capacity: 100,
		Emit: func(emitter stream.Emitter) {
			conn, err := net.Dial(protocol, addr)
			if err != nil {
				return
			}

			scanner.Attach(conn)
			for {
				if message, err := scanner.Scan(); err == nil {
					emitter.Emit(message)
					continue
				}
				return
			}
		},
	}
}
