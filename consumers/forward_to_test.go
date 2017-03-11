package consumers_test

import (
	"testing"

	"github.com/drborges/rivers/consumers"
	. "github.com/drborges/rivers/consumers/matchers"
	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/stream"
)

func TestForwardToConsumer(t *testing.T) {
	expect := expectations.New()

	ch := make(chan stream.T, 3)
	consumer := consumers.ForwardTo(ch)

	if err := expect(consumer).To(Forward(1, 2, 3).To(ch)); err != nil {
		t.Error(err)
	}
}

func TestForwardToDoesNotConsumeDataIfUpstreamIsClosed(t *testing.T) {
	expect := expectations.New()

	ch := make(chan stream.T, 0)
	consumer := consumers.ConsumerWithClosedUpstream(consumers.ForwardTo(ch))

	if err := expect(consumer).ToNot(Forward(1, 2, 3).To(ch)); err != nil {
		t.Error(err)
	}
}
