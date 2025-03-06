package util

import (
	"reflect"
	"strings"
	"time"
)

// IsEmptyValue func interface value checker
func IsEmptyValue(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(v.String())) == 0
	case reflect.Array, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Bool:
		if v.Bool() == false {
			return false
		}

		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmptyValue(v.Elem().Interface())
	case reflect.Struct:
		vl, ok := i.(time.Time)
		if ok && vl.IsZero() {
			return true
		}
	}

	return false
}
