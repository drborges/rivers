package producers

import (
	"bufio"
	"github.com/drborges/rivers/stream"
	"io"
	"os"
	"reflect"
)

func FromRange(from, to int) stream.Producer {
	return &Observable{
		Capacity: to - from + 1,
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

func FromReader(r io.Reader) stream.Producer {
	return &Observable{
		Emit: func(emitter stream.Emitter) {
			buf := bufio.NewReader(r)

			for {
				b, err := buf.ReadByte()
				if err == io.EOF {
					return
				}
				if err != nil {
					panic(err)
				}

				emitter.Emit(b)
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
