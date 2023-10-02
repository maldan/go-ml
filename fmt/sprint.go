package mfmt

import (
	mmath "github.com/maldan/go-ml/math"
	"reflect"
)

func sprintNil() string {
	if UseColor {
		return FgGray + "nil" + Reset
	} else {
		return "nil"
	}
}

func sprintString(s string) string {
	if UseColor {
		return FgGreen + s + Reset
	} else {
		return s
	}
}

func sprintBool(b bool) string {
	if UseColor {
		if b {
			return FgYellow + "true" + Reset
		} else {
			return FgYellow + "false" + Reset
		}
	} else {
		if b {
			return "true"
		} else {
			return "false"
		}
	}
}

func sprintFloatFraction(f float64) string {
	result := ""
	for i := 0; i < 10; i++ { // Преобразование 10 десятичных знаков
		f *= 10
		digit := int64(f)
		result += string(rune('0' + digit))
		f -= float64(digit)
	}
	return result
}

func sprintFloat(f float64) string {
	integerPart := int64(f)
	fractionalPart := mmath.Abs(f) - mmath.Abs(float64(integerPart))

	intStr := sprintInt(integerPart)
	fracStr := sprintFloatFraction(fractionalPart)

	// If -0.0...
	if f < 0 && integerPart == 0 {
		intStr = "-" + intStr
	}

	if fractionalPart == 0 {
		if UseColor {
			return FgRed + intStr + Reset
		} else {
			return intStr
		}
	}

	if UseColor {
		return intStr + FgRed + "." + fracStr + Reset
	} else {
		return intStr + "." + fracStr
	}
}

func sprintInt(n int64) string {
	if n == 0 {
		if UseColor {
			return FgRed + "0" + Reset
		} else {
			return "0"
		}
	}

	isNegative := false
	if n < 0 {
		isNegative = true
		n = -n
	}

	var result string

	for n > 0 {
		digit := n % 10
		result = string(rune('0'+digit)) + result
		n /= 10
	}

	if isNegative {
		result = "-" + result
	}

	if UseColor {
		return FgRed + result + Reset
	} else {
		return result
	}
}

func sprintStruct(ss any) string {
	out := "{ "
	t := reflect.TypeOf(ss)
	v := reflect.ValueOf(ss)
	for i := 0; i < t.NumField(); i++ {
		el := v.Field(i)
		/*if el.IsNil() {
			out += sprintNil()
			continue
		}*/

		elType := t.Field(i).Type.Kind()

		// out += " "
		out += t.Field(i).Name + ": "

		switch elType {
		case reflect.Struct:
			out += sprintStruct(el.Interface())
			break
		case reflect.Slice, reflect.Array:
			out += sprintSlice(el.Interface())
			break
		case reflect.Float32, reflect.Float64:
			out += sprintFloat(el.Float())
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out += sprintInt(el.Int())
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			out += sprintInt(int64(el.Uint()))
			break
		case reflect.Bool:
			out += sprintBool(el.Bool())
			break
		case reflect.String:
			out += sprintString(el.String())
			break
		default:
			out += v.Field(i).String()
			break
		}

		if i < t.NumField()-1 {
			out += ", "
		}

		// out += "\n"
	}
	return out + " }"
}

func sprintSlice(ss any) string {
	out := "["
	v := reflect.ValueOf(ss)

	for i := 0; i < v.Len(); i++ {
		el := v.Index(i)
		/*if el.IsNil() {
			out += sprintNil()
			continue
		}
		*/
		elType := el.Kind()

		switch elType {
		case reflect.Float32, reflect.Float64:
			out += sprintFloat(el.Float())
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			out += sprintInt(el.Int())
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			out += sprintInt(int64(el.Uint()))
			break
		case reflect.Struct:
			out += sprintStruct(el.Interface())
			break
		case reflect.Slice, reflect.Array:
			out += sprintSlice(el.Interface())
			break
		case reflect.Bool:
			out += sprintBool(el.Bool())
			break
		case reflect.String:
			out += sprintString(el.String())
			break
		default:
			out += el.String()
			break
		}

		if i < v.Len()-1 {
			out += ", "
		}
	}
	return out + "]"
}
