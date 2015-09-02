package where

import (
	"github.com/drborges/rivers/stream"
	"reflect"
	"regexp"
)

func ItemResemble(subject stream.T) stream.PredicateFn {
	return func(data stream.T) bool {
		return reflect.DeepEqual(data, subject)
	}
}

func ItemIs(subject stream.T) stream.PredicateFn {
	return func(data stream.T) bool {
		return data == subject
	}
}

func StructHas(field string, value stream.T) stream.PredicateFn {
	return func(data stream.T) bool {
		typ := reflect.TypeOf(data)
		val := reflect.ValueOf(data)

		if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
			return false
		}

		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
			return false
		}

		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		}

		return val.FieldByName(field).Interface() == value
	}
}

func StructFieldMatches(field string, pattern string) stream.PredicateFn {
	return func(data stream.T) bool {
		typ := reflect.TypeOf(data)
		val := reflect.ValueOf(data)

		if typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
			return false
		}

		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
			return false
		}

		if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
			val = val.Elem()
		}

		r := regexp.MustCompile(pattern)
		return r.MatchString(val.FieldByName(field).Interface().(string))
	}
}
