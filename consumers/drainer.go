package consumers

import "github.com/drborges/riversv2/rx"

type drainer struct {
	context rx.Context
}

func (drainer *drainer) Consume(in rx.InStream) {
	for range in {}
}