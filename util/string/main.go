package ml_string

import (
	"fmt"
	"strconv"
	"strings"
)

func UnTitle(str string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return strings.ToLower(str[0:1])
	}
	return strings.ToLower(str[0:1]) + str[1:]
}

func Title(str string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return strings.ToUpper(str[0:1])
	}
	return strings.ToUpper(str[0:1]) + str[1:]
}

func ToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func ToFloat(s string) float64 {
	n, _ := strconv.ParseFloat(s, 64)
	return n
}

func OnlyDigit(str string) string {
	return Only(str, "0123456789")
}

func Only(str string, allowList string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(allowList, r) {
			return r
		}
		return -1
	}, str)
}

func NonNull(str any) string {
	if str == nil {
		return ""
	}
	return fmt.Sprintf("%v", str)
}

func FromAny(a any) string {
	return fmt.Sprintf("%v", a)
}
