package stream

func New(capacity int) (Readable, Writable) {
	ch := make(chan T, capacity)
	return ch, ch
}

func NewEmpty() Readable {
	ch := make(chan T)
	close(ch)
	return ch
}

func (readable Readable) Read() []T {
	read := []T{}
	for data := range readable {
		read = append(read, data)
	}
	return read
}
