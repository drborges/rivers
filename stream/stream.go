package stream

func New(capacity int) (Readable, Writable) {
	ch := make(chan T, capacity)
	return ch, ch
}

func (readable Readable) ReadAll() []T {
	read := []T{}
	for data := range readable {
		read = append(read, data)
	}
	return read
}

func (readable Readable) Capacity() int {
	return cap(readable)
}
