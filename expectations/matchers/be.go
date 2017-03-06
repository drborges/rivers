package matchers

import (
	"errors"
	"fmt"

	"github.com/drborges/rivers/expectations"
)

// Be provides an expectations.MatchFunc that verifies whether the value under
// test holds the same reference as the expected one.
func Be(expected interface{}) expectations.MatchFunc {
	return func(actual interface{}) error {
		if actual != expected {
			return errors.New(fmt.Sprintf("Expected %v, to be %v", expected, actual))
		}

		return nil
	}
}
