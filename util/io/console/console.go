package ml_console

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
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

func PrettyPrint(x any) {
	// Create a struct describing your scheme
	scheme := pp.ColorScheme{
		//Integer:       pp.Green | pp.Bold,
		//Float:         pp.Black | pp.BackgroundWhite | pp.Bold,
		String:          pp.Green | pp.Bold,
		StringQuotation: pp.Green | pp.Bold,
		FieldName:       pp.Magenta,
		Nil:             pp.Red,
	}

	// Register it for usage
	pp.Default.SetColorScheme(scheme)

	pp.Println(x)
}
