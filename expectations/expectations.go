package expectations

import "errors"

type Expectation interface {
	To(Matcher) error
	ToNot(Matcher) error
}

type Matcher interface {
	Match(actual interface{}) error
}

type ExpectFunc func(actual interface{}) Expectation

type MatchFunc func(actual interface{}) error

func (fn MatchFunc) Match(actual interface{}) error {
	return fn(actual)
}

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
