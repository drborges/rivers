package matchers

import (
	"errors"
	"fmt"
	"time"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

func TimeoutWithin(duration time.Duration) expectations.MatchFunc {
	return func(actual interface{}) error {
		writer, ok := actual.(stream.Writer)

		if !ok {
			return errors.New(fmt.Sprintf("Expected an actual that implements stream.Writer, got %v", actual))
		}

		timeout := make(chan bool, 1)
		go func() {
			writer.Write(1)
			timeout <- true
		}()

		timeoutExceeded := make(chan bool, 1)
		go func() {
			time.Sleep(duration + 10*time.Millisecond)
			timeoutExceeded <- true
		}()

		select {
		case <-timeoutExceeded:
			return errors.New(fmt.Sprintf("Expected stream.Writer to have timed out within %v", duration))
		case <-timeout:
			return nil
		}
	}
}
