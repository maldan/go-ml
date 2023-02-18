package gosn

import (
	"encoding/binary"
	"fmt"
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
	fmt.Printf("%v ", TypeToString(tp))

	if tp == T_STRUCT {
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

	if tp == T_STRING {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		fmt.Printf("\"%v\"", string(blob))
	}

	if tp == T_BOOL {
		fmt.Printf("%v", bytes[offset] == 1)
		offset += 1
	}

	if tp == T_32 {
		value := binary.LittleEndian.Uint32(bytes[offset:])
		offset += 4
		fmt.Printf("%v", value)
	}

	if tp == T_64 {
		value := binary.LittleEndian.Uint64(bytes[offset:])
		offset += 8
		fmt.Printf("%v", value)
	}

	if tp == T_SLICE {
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

	if tp == T_CUSTOM {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		fmt.Printf("%v", blob)
	}

	return offset
}
