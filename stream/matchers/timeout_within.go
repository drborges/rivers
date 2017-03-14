package matchers

import (
	"fmt"
	"time"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

// TimeoutWithin matcher that verifies if the given stream.Writer times out within
// the given time.Duration.
func TimeoutWithin(duration time.Duration) expectations.MatchFunc {
	return func(actual interface{}) error {
		writer, ok := actual.(stream.Writer)

		if !ok {
			return fmt.Errorf("Expected an actual that implements stream.Writer, got %v", actual)
		}

		timeout := make(chan bool, 0)
		go func() {
			defer close(timeout)
			writer.Write(1)
		}()

		select {
		case <-time.After(duration + 10*time.Millisecond):
			return fmt.Errorf("Expected stream.Writer to have timed out within %v", duration)
		case <-timeout:
			return nil
		}
	}
}
