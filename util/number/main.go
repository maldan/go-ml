package ml_number

import (
	"fmt"
	"math"
	"reflect"
)

func ToLength[T any](x any) uint64 {
	typeOf := reflect.TypeOf(x)

	switch typeOf.Kind() {
	case reflect.Uint8:
		return uint64(x.(uint8))
	case reflect.Uint16:
		return uint64(x.(uint16))
	case reflect.Uint32:
		return uint64(x.(uint32))
	case reflect.Uint64:
		return x.(uint64)
	case reflect.Uint:
		return uint64(x.(uint))

	case reflect.Int8:
		return uint64(x.(int8))
	case reflect.Int16:
		return uint64(x.(int16))
	case reflect.Int32:
		return uint64(x.(int32))
	case reflect.Int64:
		return uint64(x.(int64))
	case reflect.Int:
		return uint64(x.(int))
	case reflect.Slice:
		valueOf := reflect.ValueOf(x)
		return uint64(valueOf.Len())
	case reflect.String:
		valueOf := reflect.ValueOf(x)
		return uint64(valueOf.Len())
	default:
		return 0
	}
}

func ToHumanReadableSize[T any](something T) string {
	length := ToLength[T](something)

	if float64(length) < math.Pow(2.0, 10.0) {
		return fmt.Sprintf("%d B", length)
	}
	if float64(length) < math.Pow(2.0, 20.0) {
		return fmt.Sprintf("%.2f kB", float64(length)/math.Pow(2.0, 10.0))
	}
	if float64(length) < math.Pow(2.0, 30.0) {
		return fmt.Sprintf("%.2f mB", float64(length)/math.Pow(2.0, 20.0))
	}
	if float64(length) < math.Pow(2.0, 40.0) {
		return fmt.Sprintf("%.2f gB", float64(length)/math.Pow(2.0, 30.0))
	}
	return fmt.Sprintf("%d", length)
}
