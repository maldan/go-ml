package ml_convert

import (
	"encoding/base64"
	"errors"
	"strings"
)

func DataUrlToBytes(dataUrl string) ([]byte, string, error) {
	if dataUrl[0:4] != "data" {
		return nil, "", errors.New("not a data url")
	}

	isFound := false
	slice := dataUrl[:]
	dataType := ""
	for i := 0; i < 32; i++ {
		if dataUrl[i] == ',' {
			tuple := strings.Split(dataUrl[5:i], ";")
			dataType = tuple[0]
			slice = dataUrl[i+1:]
			isFound = true
			break
		}
	}
	if !isFound {
		return nil, "", errors.New("not a data url")
	}

	decoded, err := base64.StdEncoding.DecodeString(slice)
	if err != nil {
		return nil, "", err
	}

	return decoded, dataType, nil
}
