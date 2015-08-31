package stream

func New(capacity int) (Readable, Writable) {
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
