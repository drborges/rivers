package producers_test

import (
	"testing"

	"github.com/drborges/rivers/expectations"
	"github.com/drborges/rivers/producers"
	. "github.com/drborges/rivers/producers/matchers"
)

func TestRangeProducesData(t *testing.T) {
	expect := expectations.New()

	producer := producers.Range(0, 3)

	if err := expect(producer).To(Produce(0, 1, 2, 3)); err != nil {
		t.Error(err)
	}
}

func TestRangeDoesNotProduceWhenContextIsClosed(t *testing.T) {
	expect := expectations.New()

	producer := producers.WithClosedContext(producers.Range(0, 3))

	if err := expect(producer).ToNot(Produce(0, 1, 2, 3)); err != nil {
		t.Error(err)
	}
}
