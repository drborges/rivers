package stream

import "context"

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

	// Fork creates a new pair of stream reader and writer
	Fork() (Reader, Writer)
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
	ch := make(chan T, 2)
	rootCtx, cancelRoot := context.WithCancel(context.Background())
	ctx, cancel := contextFromRoot(rootCtx)

	return &reader{ctx, cancel, rootCtx, cancelRoot, ch}, &writer{ctx, cancelRoot, ch}
}

func contextFromRoot(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

type reader struct {
	ctx        context.Context
	cancel     context.CancelFunc
	rootCtx    context.Context
	cancelRoot context.CancelFunc
	ch         Readable
}

func (reader *reader) Read() Readable {
	return reader.ch
}

func (parent *reader) Fork() (Reader, Writer) {
	ch := make(chan T, 2)
	ctx, cancel := contextFromRoot(parent.rootCtx)
	propagateCancellation := context.CancelFunc(func() {
		cancel()
		parent.cancel()
	})

	return &reader{ctx, propagateCancellation, parent.rootCtx, parent.cancelRoot, ch}, &writer{ctx, parent.cancelRoot, ch}
}

func (reader *reader) Close(err error) {
	reader.cancel()
	if err != nil {
		reader.cancelRoot()
	}
}

type writer struct {
	ctx        context.Context
	cancelRoot context.CancelFunc
	ch         Writable
}

func (writer *writer) Write(data T) error {
	writer.ch <- data
	return nil
}

func (writer *writer) Close(err error) {
	defer close(writer.ch)
	if err != nil {
		writer.cancelRoot()
	}
}
