package ml_reflect

import (
	"fmt"
	"reflect"
)

func SetStructField(s any, name string, value any) bool {
	// Nil check
	if s == nil {
		return false
	}

	// Pointer check
	typeOf := reflect.TypeOf(s)
	if typeOf.Kind() != reflect.Pointer {
		return false
	}

	// Pointer on struct
	if typeOf.Elem().Kind() != reflect.Struct {
		return false
	}

	// Get value
	valueOf := reflect.ValueOf(s)
	field := valueOf.Elem().FieldByName(name)
	if field.IsZero() {
		return false
	}
	if !field.IsValid() {
		return false
	}
	if !field.CanSet() {
		return false
	}

	if field.Kind() == reflect.Pointer {

	} else {
		v := reflect.ValueOf(value)
		fmt.Printf("%v\n", field.Type())
		field.Set(v)
	}

	return true
}
