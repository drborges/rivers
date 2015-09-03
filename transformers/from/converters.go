package from

import (
	"encoding/json"
	"github.com/drborges/rivers/stream"
	"reflect"
	"fmt"
)

func StructToJSON(data stream.T) stream.T {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

func JSONToStruct(example interface{}) stream.MapFn {
	typ := reflect.TypeOf(example)

	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() != reflect.Struct ||
		typ.Kind() != reflect.Ptr && typ.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Expected a struct or a struct pointer, got %v", typ.Kind()))
	}

	dst := reflect.New(typ).Interface()
	return func(data stream.T) stream.T {
		if err := json.Unmarshal(data.([]byte), dst); err != nil {
			panic(err)
		}
		return reflect.ValueOf(dst).Elem().Interface()
	}
}
