package where

import (
	"github.com/drborges/rivers/stream"
	"reflect"
	"regexp"
	"strings"
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

		// Handle inner structs/arrays/slices
		if strings.Contains(field, ".") {
			path := strings.Split(field, ".")
			pathHead := path[0]
			pathTail := strings.Join(path[1:], ".")
			currentField := val.FieldByName(pathHead)

			if !currentField.IsValid() {
				return false
			}

			// Handle field slice/array
			if currentField.Kind() == reflect.Array || currentField.Kind() == reflect.Slice {
				for i := 0; i < currentField.Len(); i++ {
					if StructHas(pathTail, value)(currentField.Index(i).Interface()) {
						return true
					}
				}
				return false
			}

			fieldData := currentField.Interface()
			return StructHas(pathTail, value)(fieldData)
		}

		fieldVal := val.FieldByName(field)

		if !fieldVal.IsValid() {
			return false
		}

		// Handle field slice/array
		if fieldVal.Kind() == reflect.Array || fieldVal.Kind() == reflect.Slice {
			for i := 0; i < fieldVal.Len(); i++ {
				if fieldVal.Index(i).Interface() == value {
					return true
				}
			}
			return false
		}

		return fieldVal.Interface() == value
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
