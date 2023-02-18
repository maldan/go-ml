package ml_gosn_test

import (
	"github.com/maldan/go-ml/util/encode/gosn"
	ml_time "github.com/maldan/go-ml/util/time"
	"testing"
)

type TestStruct struct {
	Bool   bool
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64

	Float32 float32
	Float64 float64

	String string
	Slice  []string
	Struct struct {
		A uint8
		B string
	}
	Time ml_time.DateTime
}

type TestStructSliceOnly struct {
	SliceStr []string
	SliceU8  []uint8
}

func sas(t *testing.T, nameToId ml_gosn.NameToId) {
	tsIn := TestStruct{
		Bool:    true,
		Uint8:   255,
		Uint16:  65535,
		Uint32:  100,
		Uint64:  200,
		Float32: 0.0001,
		Float64: 0.000001,
		String:  "ABABAB",
		Slice:   []string{"A", "B", "C"},
		Struct: struct {
			A uint8
			B string
		}{A: 32, B: "B"},
		Time: ml_time.Now(),
	}

	// Prepare
	x := make([]byte, 0)
	tsOut := TestStruct{}

	// Encode and decode
	if nameToId == nil {
		x = ml_gosn.Marshal(tsIn)
		ml_gosn.Unmarshall(x, &tsOut)
	} else {
		x = ml_gosn.MarshalExt(tsIn, nameToId)
		ml_gosn.UnmarshallExt(x, &tsOut, nameToId.Invert())
	}

	// Compare bool
	if tsIn.Bool != tsOut.Bool {
		t.Fatal("fuck")
	}

	// Compare numbers
	if tsIn.Uint8 != tsOut.Uint8 {
		t.Fatal("fuck")
	}
	if tsIn.Uint16 != tsOut.Uint16 {
		t.Fatal("fuck")
	}
	if tsIn.Uint32 != tsOut.Uint32 {
		t.Fatal("fuck")
	}
	if tsIn.Uint64 != tsOut.Uint64 {
		t.Fatal("fuck")
	}

	// Compare floats
	if tsIn.Float32 != tsOut.Float32 {
		t.Fatal("fuck")
	}
	if tsIn.Float64 != tsOut.Float64 {
		t.Fatal("fuck")
	}

	// Compare string
	if tsIn.String != tsOut.String {
		t.Fatal("fuck")
	}

	// Compare slice
	if len(tsIn.Slice) != len(tsOut.Slice) {
		t.Fatal("fuck")
	}
	for i := 0; i < len(tsIn.Slice); i++ {
		if tsIn.Slice[i] != tsOut.Slice[i] {
			t.Fatalf("fuck %v - %v", tsIn.Slice[i], tsOut.Slice[i])
		}
	}

	// Compare struct
	if tsIn.Struct != tsOut.Struct {
		t.Fatal("fuck")
	}

	// Compare custom
	if tsIn.Time != tsOut.Time {
		t.Fatal("fuck")
	}
}

func TestMain_Named(t *testing.T) {
	sas(t, nil)
}

func TestMain_Id(t *testing.T) {
	nameToId := ml_gosn.NameToId{}
	nameToId.FromStruct(TestStruct{})
	sas(t, nameToId)
}

func TestMainSlice(t *testing.T) {
	tsIn := TestStructSliceOnly{
		SliceStr: []string{"A", "B", "C"},
		SliceU8:  []uint8{1, 2, 3},
	}

	// Encode
	x := ml_gosn.Marshal(tsIn)

	// Decode
	tsOut := TestStructSliceOnly{}
	ml_gosn.Unmarshall(x, &tsOut)

	// Compare slice
	if len(tsIn.SliceStr) != len(tsOut.SliceStr) {
		t.Fatal("fuck")
	}

	for i := 0; i < len(tsIn.SliceStr); i++ {
		if tsIn.SliceStr[i] != tsOut.SliceStr[i] {
			t.Fatalf("fuck %v - %v", tsIn.SliceStr[i], tsOut.SliceStr[i])
		}
	}
}
