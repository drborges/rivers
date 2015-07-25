package scanners

import (
	"bufio"
	"io"
)

type BufioScanner struct {
	scanner  *bufio.Scanner
	splitter bufio.SplitFunc
}

func NewWordScanner() Scanner {
	return &BufioScanner{
		splitter: bufio.ScanWords,
	}
}

func NewLineScanner() Scanner {
	return &BufioScanner{
		splitter: bufio.ScanLines,
	}
}

func (bufioscan *BufioScanner) Attach(r io.Reader) {
	bufioscan.scanner = bufio.NewScanner(r)
}

func (bufioscan *BufioScanner) Scan() ([]byte, error) {
	bufioscan.scanner.Split(bufioscan.splitter)
	if ok := bufioscan.scanner.Scan(); ok {
		return []byte(bufioscan.scanner.Text()), nil
	}
	return []byte{}, io.EOF
}
