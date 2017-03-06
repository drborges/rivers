package expectations

import "errors"

// Expectation represents a test expectation with either a positive or negative
// outcome.
type Expectation interface {
	// To applies the given matcher to the value under test, returning nil if the
	// value matches the expectation. An error is returned otherwise.
	To(Matcher) error
	// ToNot applies the given matcher to the value under test, returning nil if the
	// value does not match the expectation. An error is returned otherwise.
	ToNot(Matcher) error
}

// Matcher implements the expectation condition
type Matcher interface {
	// Match checks the value under test for a given condition. An error is
	// returned in case the value does not match the condition, otherwise, nil is
	// returned.
	Match(actual interface{}) error
}

// ExpectFunc function that takes the value under test and returns the expectation
// chainning API for applying the matchers.
type ExpectFunc func(actual interface{}) Expectation

// MatchFunc shorthand function that implements the Matcher interface.
type MatchFunc func(actual interface{}) error

// Match checks the value under test against a matching condition. Returns an
// error if the value does not match the condition.
func (fn MatchFunc) Match(actual interface{}) error {
	return fn(actual)
}

// New creates a new expectation function for testing a given value against
// possible matching conditions.
func New() ExpectFunc {
	return func(actual interface{}) Expectation {
		return &expectation{actual: actual}
	}
}

type expectation struct {
	actual interface{}
}

func (expect *expectation) To(matcher Matcher) error {
	return matcher.Match(expect.actual)
}

func (expect *expectation) ToNot(matcher Matcher) error {
	if err := matcher.Match(expect.actual); err == nil {
		return errors.New("Expected actual to have not matched expectation")
	}

	return nil
}
