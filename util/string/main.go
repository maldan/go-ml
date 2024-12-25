package ml_string

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const ENGLISH_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const DIGITS = "0123456789"

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

func OnlyUID(str string) string {
	return Only(str, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
}

func Only(str string, charset string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune(charset, r) {
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

func PascalCaseToSnakeCase(s string) string {
	re := regexp.MustCompile("([A-Z][a-z0-9]+)")
	words := re.FindAllString(s, -1)

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "_")
}
