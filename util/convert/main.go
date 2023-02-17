package ml_convert

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func DataUrlToBytes(dataUrl string) ([]byte, string, error) {
	if dataUrl[0:4] != "data" {
		return nil, "", errors.New("data segment not found, not a data url")
	}

	isFound := false
	slice := dataUrl[:]
	dataType := ""
	for i := 0; i < len(dataUrl); i++ {
		if dataUrl[i] == ',' {
			tuple := strings.Split(dataUrl[5:i], ";")
			dataType = tuple[0]
			slice = dataUrl[i+1:]
			isFound = true
			break
		}
	}
	if !isFound {
		return nil, "", errors.New("comma not found, not a data url")
	}

	decoded, err := base64.StdEncoding.DecodeString(slice)
	if err != nil {
		return nil, "", err
	}

	return decoded, dataType, nil
}

func ToBase64[T string | []byte](v T) string {
	switch any(v).(type) {
	case string:
		enc := base64.URLEncoding.EncodeToString([]byte(any(v).(string)))
		return enc
	case []byte:
		enc := base64.URLEncoding.EncodeToString(any(v).([]byte))
		return enc
	default:
		return ""
	}
}

func ToString[T any](v T) string {
	v2 := any(v)
	switch v2.(type) {
	case int:
		return strconv.Itoa(v2.(int))
	case uint8:
		return strconv.Itoa(int(v2.(uint8)))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func FromBase64(v string) []byte {
	uDec, _ := base64.URLEncoding.DecodeString(v)
	return uDec
}
