package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-ml/db/goson/core"
)

func visualizeIntent(amount int) {
	for i := 0; i < amount; i++ {
		fmt.Printf(" ")
	}
}

func Visualize(bytes []byte, intent int) int {
	offset := 0

	// Read type
	tp := bytes[offset]
	offset += 1
	fmt.Printf("%v ", core.TypeToString(tp))

	if tp == core.TypeStruct {
		fmt.Printf("\n")
		visualizeIntent(intent)
		fmt.Printf("{\n")

		// Size
		size := binary.LittleEndian.Uint16(bytes[offset:])
		offset += 2
		visualizeIntent(intent + 2)
		fmt.Printf("Size: %v\n", size)

		// Amount
		amount := int(bytes[offset])
		offset += 1
		visualizeIntent(intent + 2)
		fmt.Printf("Fields: %v\n\n", amount)

		// Go over fields
		for i := 0; i < amount; i++ {
			fieldId := bytes[offset]
			visualizeIntent(intent + 2)
			fmt.Printf("%v: ", fieldId)
			offset += 1

			offset += Visualize(bytes[offset:], intent+2)
			fmt.Printf("\n")
		}

		visualizeIntent(intent)
		fmt.Printf("}\n")
	}

	if tp == core.TypeString {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		fmt.Printf("\"%v\"", string(blob))
	}

	if tp == core.T_BOOL {
		fmt.Printf("%v", bytes[offset] == 1)
		offset += 1
	}

	if tp == core.T_32 {
		value := binary.LittleEndian.Uint32(bytes[offset:])
		offset += 4
		fmt.Printf("%v", value)
	}

	if tp == core.T_64 {
		value := binary.LittleEndian.Uint64(bytes[offset:])
		offset += 8
		fmt.Printf("%v", value)
	}

	if tp == core.TypeSlice {
		fmt.Printf("[\n")

		// Size
		size := binary.LittleEndian.Uint16(bytes[offset:])
		offset += 2
		visualizeIntent(intent + 2)
		fmt.Printf("Size: %v\n", size)

		// Amount
		amount := int(bytes[offset])
		offset += 2
		visualizeIntent(intent + 2)
		fmt.Printf("Fields: %v\n\n", amount)

		for i := 0; i < amount; i++ {
			offset += Visualize(bytes[offset:], intent+2)
		}

		visualizeIntent(intent)
		fmt.Printf("]\n")
	}

	if tp == core.T_CUSTOM {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		fmt.Printf("%v", blob)
	}

	/*if tp == core.TypeTime {
		size := int(bytes[offset])
		offset += 1
		blob := bytes[offset : offset+size]

		x, err := time.Parse("2006-01-02T15:04:05.999-07:00", string(blob))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("%v", x)

		offset += size
	}*/

	return offset
}
