package mfmt

import (
	"reflect"
)

var UseColor = false

func Sprintf(format string, a ...any) string {
	out := make([]byte, 0, len(format))
	argId := 0

	if UseColor {
		out = append(out, []byte(Reset)...)
	}

	for i := 0; i < len(format)-1; i++ {
		if format[i] == '%' && format[i+1] == 'v' {
			if a[argId] == nil {
				s := sprintNil()
				out = append(out, []byte(s)...)
				continue
			}
			argType := reflect.TypeOf(a[argId]).Kind()

			if argType == reflect.Int {
				s := sprintInt(int64(a[argId].(int)))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Int8 {
				s := sprintInt(int64(a[argId].(int8)))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Int16 {
				s := sprintInt(int64(a[argId].(int16)))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Int32 {
				s := sprintInt(int64(a[argId].(int32)))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Int64 {
				s := sprintInt(a[argId].(int64))
				out = append(out, []byte(s)...)
			}

			if argType == reflect.Float32 {
				s := sprintFloat(float64(a[argId].(float32)))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Float64 {
				s := sprintFloat(a[argId].(float64))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Struct {
				s := sprintStruct(a[argId])
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Slice || argType == reflect.Array {
				s := sprintSlice(a[argId])
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Bool {
				s := sprintBool(a[argId].(bool))
				out = append(out, []byte(s)...)
			}
			if argType == reflect.String {
				s := sprintString(a[argId].(string))
				out = append(out, []byte(s)...)
			}

			argId += 1
			i += 1
		} else {
			out = append(out, format[i])
		}
	}

	// Last char
	// out = append(out, format[len(format)-1])

	return string(out)
}
