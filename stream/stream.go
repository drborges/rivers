package stream

import (
	goContext "context"

	"github.com/drborges/rivers/context"
)

// T data type flowing through rivers streams
type T interface{}

// Readable a channel which one may read data from
type Readable <-chan T

// Writable a channel which one may write data to
type Writable chan<- T

// Closer represents a type which may be closed. When no error is given, that
// indicates the Closer was gracefully closed without errors. Closing with an
// error, indicates a failure.
type Closer interface {
	Close(error)
}

// Reader provides means to read from a readable stream as well as signal the
// termination of the stream, gracefully or not.
type Reader interface {
	Closer

	// Read provides a readable stream from which data can be read
	Read() Readable
}

// Writer provides means to write data to a writable stream as well as signal the
// termination of the stream, gracefully or not. Closing also closes the underlying
// writable stream.
type Writer interface {
	Closer

	// Write writes the given data to the underlying writable stream, returning an
	// error in case of a failure.
	Write(data T) error
}

// Empty represents an empty readable stream which has already ceased producing
// data.
var Empty = func() Reader {
	r, w := New()
	w.Close(nil)
	return r
}()

// New Creates the Reader and Writer components of a rivers stream with the default
// configuration.
func New() (Reader, Writer) {
	return NewWithStdContext(goContext.Background())
}

// NewWithStdContext Creates the Reader and Writer components of a rivers stream
// with the default configuration.
func NewWithStdContext(stdCtx goContext.Context) (Reader, Writer) {
	ch := make(chan T, 2)
	ctx := context.FromStdContext(stdCtx)
	return &reader{ctx, ch}, &writer{ctx, ch}
}

type reader struct {
	ctx context.Context
	ch  Readable
}

func (reader *reader) Read() Readable {
	return reader.ch
}

func (reader *reader) Close(err error) {
	reader.ctx.Close()
}

type writer struct {
	ctx context.Context
	ch  Writable
}

func (writer *writer) Write(data T) error {
	select {
	case writer.ch <- data:
	}
	return nil
}

func (writer *writer) Close(err error) {
	defer close(writer.ch)
	writer.ctx.Close()
}
