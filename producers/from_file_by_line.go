package producers

import (
	"github.com/drborges/riversv2/rx"
	"os"
	"bufio"
)

type fromFileByLine struct {
	context  rx.Context
	file *os.File
}

func (p *fromFileByLine) Produce() rx.InStream {
	// TODO find a better way to set stream capacity
	reader, writer := rx.NewStream(100)

	go func() {
		defer p.context.Recover()
		defer p.file.Close()
		defer close(writer)

		scanner := bufio.NewScanner(p.file)
		for scanner.Scan() {
			select {
			case <-p.context.Closed():
				return
			default:
				writer <- scanner.Text()
			// In case of error the context should be closed
			// if scanner.Err() != nil {
			//     p.context.Close(scanner.Err())
			// }
			}
		}
	}()

	return reader
}
