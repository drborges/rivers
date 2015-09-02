package where

import (
	"github.com/drborges/rivers/stream"
	"reflect"
)

func ItemResemble(subject stream.T) stream.PredicateFn {
	return func(data stream.T) bool {
		return reflect.DeepEqual(data, subject)
	}
}
