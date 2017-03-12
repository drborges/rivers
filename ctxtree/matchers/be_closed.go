package matchers

import (
	"errors"
	"fmt"

	"github.com/drborges/rivers/ctxtree"
	"github.com/drborges/rivers/expectations"
)

// BeClosed matcher that verifies if the given context.Context is closed.
func BeClosed() expectations.MatchFunc {
	return func(actual interface{}) error {
		ctx, ok := actual.(ctxtree.Signaler)

		if !ok {
			return fmt.Errorf("Exected an actual that implements 'ctxtree.Signaler', got %v", actual)
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			return errors.New("Expected context to be closed, it was not")
		}
	}
}
