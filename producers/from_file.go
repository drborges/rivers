package producers

import (
	"bufio"
	"github.com/drborges/rivers/stream"
	"os"
	"strings"
)

type fromFile struct {
	file *os.File
}

func (builder *fromFile) ByLine() stream.Producer {
	return &Observable{
		Capacity: 100,
		Emit: func(emitter stream.Emitter) {
			defer builder.file.Close()
			scanner := bufio.NewScanner(builder.file)
			for scanner.Scan() {
				emitter.Emit(scanner.Text())
			}
		},
	}
}

func (builder *fromFile) ByDelimiter(delimiter byte) stream.Producer {
	return &Observable{
		Capacity: 100,
		Emit: func(emitter stream.Emitter) {
			defer builder.file.Close()
			reader := bufio.NewReader(builder.file)
			for {
				line, err := reader.ReadString(delimiter)

				if line != "" {
					emitter.Emit(strings.TrimSpace(line))
				}

				if err != nil {
					return
				}
			}
		},
	}
}
