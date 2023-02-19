package ml_gosn_test

import (
	"bytes"
	"fmt"
	"github.com/maldan/go-ml/util/encode/gosn"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_time "github.com/maldan/go-ml/util/time"
	"testing"
	"text/template"
	"time"
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

func TestMainByteSlice(t *testing.T) {
	tsIn := TestStructSliceOnly{
		SliceStr: make([]string, 12),
		SliceU8:  make([]byte, 1_000_000),
	}
	for i := 0; i < len(tsIn.SliceU8); i++ {
		tsIn.SliceU8[i] = uint8(i)
	}

	// Encode
	tt := time.Now()
	x := ml_gosn.Marshal(tsIn)
	fmt.Printf("ENC: %v\n", time.Since(tt))

	// Decode
	tsOut := TestStructSliceOnly{}
	tt = time.Now()
	ml_gosn.Unmarshall(x, &tsOut)
	fmt.Printf("DEC: %v\n", time.Since(tt))

	// Compare slice
	if len(tsIn.SliceU8) != len(tsOut.SliceU8) {
		t.Fatalf("fuck %v - %v", len(tsIn.SliceU8), len(tsOut.SliceU8))
	}

	for i := 0; i < len(tsIn.SliceU8); i++ {
		if tsIn.SliceU8[i] != tsOut.SliceU8[i] {
			t.Fatalf("fuck %v - %v", tsIn.SliceU8[i], tsOut.SliceU8[i])
		}
	}
}

func _TestGen(t *testing.T) {
	tpl := `
func Test_{{ .Type }}(t *testing.T) {
	type tStruct struct {
		A {{ .Type }}
	}
	
	// Prepare
	tsIn := tStruct{ A: {{ .Value }} }
	tsOut := tStruct{}
	
	// Encode and decode
	{{ if .UserNameToId }}
    x := ml_gosn.MarshalExt(tsIn, nameToId)
	ml_gosn.UnmarshallExt(x, &tsOut, nameToId.Invert())	

	{{ else }}
	x := ml_gosn.Marshal(tsIn)
	ml_gosn.Unmarshall(x, &tsOut)
	{{ end }}
	
	// Compare
	if tsIn.A != tsOut.A {
		t.Fatalf("Not equal %v %v", tsIn.A, tsOut.A)
	}
}
	`

	// Prepare
	tt, err := template.New("code").Parse(tpl)
	if err != nil {
		panic(err)
	}

	// Execute
	var code bytes.Buffer
	code.Write([]byte("package ml_gosn_test\n"))

	typeList := []string{
		"uint8", "uint16", "uint32", "uint64",
		"int8", "int16", "int32", "int64",
		"float32", "float64", "string",
	}
	defaultValues := []any{
		255, 65535, 10241234, 3434343434343,
		-32, -12444, -10241234, -3434343434343,
		32.21222, 32.3333333344542, "\"sasageooo\"",
	}

	// Type list
	for i, tp := range typeList {
		err = tt.Execute(&code, map[string]any{
			"Type":         tp,
			"UserNameToId": false,
			"Value":        defaultValues[i],
		})
		if err != nil {
			panic(err)
		}
	}

	ml_file.New("gen_test.go").Write(code.Bytes())
}
