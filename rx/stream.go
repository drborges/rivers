package rx

import (
	"errors"
)

var (
	ErrNoSuchPointer      = errors.New("Element is not a pointer")
	ErrNoSuchSlicePointer = errors.New("Element is not a pointer to a slice")
)

func NewStream(capacity int) (Readable, Writable) {
	pipe := make(chan T, capacity)
	return pipe, pipe
}

func (readable Readable) Read() []T {
	read := []T{}
	for data := range readable {
		read = append(read, data)
	}
	return read
}
