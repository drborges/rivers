package matchers

import (
	"errors"
	"fmt"
	"time"

	"github.com/drborges/rivers/expectations"
)

// Signaler allows checking if the underlying context is done.
type Signaler interface {
	Done() <-chan struct{}
}

// BeClosed matcher that verifies if the given context.Context is closed.
func BeClosed() expectations.MatchFunc {
	return func(actual interface{}) error {
		ctx, ok := actual.(Signaler)

		if !ok {
			return fmt.Errorf("Exected an actual that provides a Done() signal method, got %v", actual)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(20 * time.Millisecond):
			return errors.New("Expected context to be closed, it was not")
		}
	}
}
