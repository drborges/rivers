package producers

import (
	"github.com/drborges/rivers/stream"
	"os"
	"bufio"
	"strings"
)

type fromFile struct {
	context stream.Context
	file    *os.File
}

func (builder *fromFile) ByLine() stream.Producer {
	return &Observable{
		Context: builder.context,
		Capacity: 100,
		Emit: func(w stream.Writable) {
			defer builder.file.Close()
			scanner := bufio.NewScanner(builder.file)
			for scanner.Scan() {
				select {
				case <-builder.context.Closed():
					return
				default:
					w <- scanner.Text()
				}
			}
		},
	}
}

func (builder *fromFile) ByDelimiter(delimiter byte) stream.Producer {
	return &Observable{
		Context: builder.context,
		Capacity: 100,
		Emit: func(w stream.Writable) {
			defer builder.file.Close()
			reader := bufio.NewReader(builder.file)
			for {
				select {
				case <-builder.context.Closed():
					return
				default:
					line, err := reader.ReadString(delimiter)

					if line != "" {
						w <- strings.TrimSpace(line)
					}

					if err != nil {
						return
					}
				}
			}
		},
	}
}
