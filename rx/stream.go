package rx

import (
	"reflect"
	"errors"
)

var (
	ErrNoSuchPointer      = errors.New("Element is not a pointer")
	ErrNoSuchSlicePointer = errors.New("Element is not a pointer to a slice")
)

func NewStream(capacity int) (InStream, OutStream) {
	pipe := make(chan T, capacity)
	return pipe, pipe
}

func (stream InStream) Read() []T {
	read := []T{}
	for data := range stream {
		read = append(read, data)
	}
	return read
}

func (stream InStream) ReadFirst() T {
	return <-stream
}

func (stream InStream) ReadLast() T {
	var last T
	for data := range stream {
		last = data
	}
	return last
}

func (stream InStream) ReadLastAs(dst interface{}) error {
	dstValue := reflect.ValueOf(dst)

	if dstValue.Kind() != reflect.Ptr {
		return ErrNoSuchPointer
	}

	data := stream.ReadLast()
	dataValue := reflect.ValueOf(data)
	dstValue.Elem().Set(dataValue)
	return nil
}

func (stream InStream) ReadFirstAs(dst interface{}) error {
	dstValue := reflect.ValueOf(dst)

	if dstValue.Kind() != reflect.Ptr {
		return ErrNoSuchPointer
	}

	data := stream.ReadFirst()
	dataValue := reflect.ValueOf(data)
	dstValue.Elem().Set(dataValue)
	return nil
}

func (stream InStream) ReadAs(dst interface{}) error {
	dstValue := reflect.ValueOf(dst)

	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Slice {
		return ErrNoSuchSlicePointer
	}

	slice := dstValue.Elem()

	for data := range stream {
		slice.Set(reflect.Append(slice, reflect.ValueOf(data)))
	}

	return nil
}

func (reader InStream) Drain() {
	for {
		_, more := <-reader
		if !more {
			return
		}
	}
}
