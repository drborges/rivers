package producers

import (
	"github.com/drborges/riversv2/rx"
	"os"
	"bufio"
	"strings"
)

type fromFileByDelimiter struct {
	context   rx.Context
	file      *os.File
	delimiter byte
}

func (p *fromFileByDelimiter) Produce() rx.InStream {
	// TODO find a better way to set stream capacity
	reader, writer := rx.NewStream(100)

	go func() {
		defer p.context.Recover()
		defer p.file.Close()
		defer close(writer)

		reader := bufio.NewReader(p.file)
		for {
			select {
			case <-p.context.Closed():
				return
			default:
				line, err := reader.ReadString(p.delimiter)

				if line != "" {
					writer <- strings.TrimSpace(line)
				}

				if err != nil {
					return
				}
			}
			// In case of error the context should be closed
			// if err != nil {
			//     p.context.Close(err)
			// }
		}
	}()

	return reader
}
