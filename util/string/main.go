package ml_string

import (
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
