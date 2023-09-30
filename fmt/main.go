package mfmt

import "reflect"

var UseColor = false

func Sprintf(format string, a ...any) string {
	out := make([]byte, 0, len(format))
	argId := 0

	if UseColor {
		out = append(out, []byte(Reset)...)
	}

	for i := 0; i < len(format)-1; i++ {
		if format[i] == '%' && format[i+1] == 'v' {
			argType := reflect.TypeOf(a[argId]).Kind()

			if argType == reflect.Int {
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
			if argType == reflect.Slice {
				s := sprintSlice(a[argId])
				out = append(out, []byte(s)...)
			}
			if argType == reflect.Bool {
				s := sprintBool(a[argId].(bool))
				out = append(out, []byte(s)...)
			}

			argId += 1
			i += 1
		} else {
			out = append(out, format[i])
		}
	}

	return string(out)
}
