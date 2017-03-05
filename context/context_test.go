package context_test

import (
	"testing"

	"github.com/drborges/rivers/context"
	. "github.com/drborges/rivers/context/matchers"
	"github.com/drborges/rivers/expectations"
)

func TestOpen(t *testing.T) {
	expect := expectations.New()

	ctx := context.New()

	if err := expect(ctx).ToNot(BeOpened()); err != nil {
		t.Error(err)
	}

	ctx.Open()

	if err := expect(ctx).To(BeOpened()); err != nil {
		t.Error(err)
	}
}

func TestOpenPropagatesSignalToChildren(t *testing.T) {
	expect := expectations.New()

	parent := context.New()
	child1 := parent.NewChild()
	child2 := parent.NewChild()

	if err := expect(parent).ToNot(BeOpened()); err != nil {
		t.Error(err)
	}

	if err := expect(child1).ToNot(BeOpened()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).ToNot(BeOpened()); err != nil {
		t.Error(err)
	}

	parent.Open()

	if err := expect(parent).To(BeOpened()); err != nil {
		t.Error(err)
	}

	if err := expect(child1).To(BeOpened()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).To(BeOpened()); err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	expect := expectations.New()

	ctx := context.New()

	if err := expect(ctx).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	ctx.Close()

	if err := expect(ctx).To(BeClosed()); err != nil {
		t.Error(err)
	}
}

func TestParentContextCannotBeClosedWhenChildrenAreStillOpened(t *testing.T) {
	expect := expectations.New()

	parent := context.New()
	child1 := parent.NewChild()
	child2 := parent.NewChild()

	if err := expect(child1).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	parent.Close()

	if err := expect(child1).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}
}

func TestParentContextIsClosedWhenAllChildrenAreClosed(t *testing.T) {
	expect := expectations.New()

	parent := context.New()
	child1 := parent.NewChild()
	child2 := parent.NewChild()
	grandchild := child2.NewChild()

	child1.Close()

	if err := expect(child1).To(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(parent).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(grandchild).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	child2.Close()

	if err := expect(child2).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(parent).ToNot(BeClosed()); err != nil {
		t.Error(err)
	}

	grandchild.Close()

	if err := expect(grandchild).To(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(child2).To(BeClosed()); err != nil {
		t.Error(err)
	}

	if err := expect(parent).To(BeClosed()); err != nil {
		t.Error(err)
	}
}
