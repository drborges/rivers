package matchers

import (
	"context"
	"errors"
	"fmt"

	"github.com/drborges/rivers/expectations"
)

func BeClosed() expectations.MatchFunc {
	return func(actual interface{}) error {
		ctx, ok := actual.(context.Context)

		if !ok {
			return errors.New(fmt.Sprintf("Exected an actual that implements 'context.Context', got %v", actual))
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			return errors.New("Expected context to be closed, it was not")
		}
	}
}
