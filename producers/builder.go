package producers

import (
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/scanners"
	"os"
	"reflect"
	"net"
)

type Builder struct {
	context stream.Context
}

func New(c stream.Context) *Builder {
	return &Builder{c}
}

func (b *Builder) FromRange(from, to int) stream.Producer {
	return &Observable{
		Context: b.context,
		Capacity: to - from,
		Emit: func(w stream.Writable) {
			for i := from; i <= to; i++ {
				select {
				case <-b.context.Closed():
					return
				default:
					w <- i
				}
			}
		},
	}
}

func (b *Builder) FromSlice(slice stream.T) stream.Producer {
	sv := reflect.ValueOf(slice)

	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Ptr {
		panic("No such slice")
	}

	if sv.Kind() == reflect.Ptr && sv.Elem().Kind() != reflect.Slice {
		panic("No such slice")
	}

	return &Observable{
		Context: b.context,
		Capacity: sv.Len(),
		Emit: func(w stream.Writable) {
			for i := 0; i < sv.Len(); i++ {
				select {
				case <-b.context.Closed():
					break
				default:
					w <- sv.Index(i).Interface()
				}
			}
		},
	}
}

func (b *Builder) FromData(data ...stream.T) stream.Producer {
	return b.FromSlice(data)
}

func (b *Builder) FromFile(f *os.File) *fromFile {
	return &fromFile{b.context, f}
}

func (b *Builder) FromSocket(protocol, addr string, scanner scanners.Scanner) stream.Producer {
	return &Observable{
		Context: b.context,
		Capacity: 100,
		Emit: func(w stream.Writable) {
			conn, err := net.Dial(protocol, addr)
			if err != nil {
				return
			}

			scanner.Attach(conn)
			for {
				select {
				case <-b.context.Closed():
					return
				default:
					if message, err := scanner.Scan(); err == nil {
						w <- message
						continue
					}
					return
				}
			}
		},
	}
}
