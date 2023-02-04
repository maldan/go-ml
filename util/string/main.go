package ml_string

import "strings"

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
