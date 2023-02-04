package ml_console

import (
	"fmt"
	"strings"
)

func PrintBytes(b []byte, lineSize int) {
	buff := make([]string, 0)
	for i := 0; i < len(b); i++ {
		buff = append(buff, fmt.Sprintf("%02X ", b[i]))
		if i != 0 && (i+1)%lineSize == 0 {
			buff = append(buff, "\n")
		}
	}
	buff = append(buff, "\n")
	fmt.Print(strings.Join(buff, ""))
}
