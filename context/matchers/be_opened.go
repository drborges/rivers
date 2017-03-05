package matchers

import (
	"errors"
	"fmt"

	"github.com/drborges/rivers/context"
	"github.com/drborges/rivers/expectations"
)

func BeOpened() expectations.MatchFunc {
	return func(actual interface{}) error {
		ctx, ok := actual.(context.Context)

		if !ok {
			return errors.New(fmt.Sprintf("Exected an actual that implements 'context.Context', got %v", actual))
		}

		select {
		case <-ctx.Opened():
			return nil
		default:
			return errors.New("Expected context to be opened, it was not")
		}
	}
}
