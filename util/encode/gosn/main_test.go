package ml_gosn_test

import (
	"bytes"
	"fmt"
	"github.com/maldan/go-ml/util/encode/gosn"
	ml_console "github.com/maldan/go-ml/util/io/console"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_time "github.com/maldan/go-ml/util/time"
	"strings"
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

/*func sas(t *testing.T, nameToId ml_gosn.NameToId) {
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
*/

/*func TestMain_Named(t *testing.T) {
	sas(t, nil)
}*/

/*func TestMain_Id(t *testing.T) {
	nameToId := ml_gosn.NameToId{}
	nameToId.FromStruct(TestStruct{})
	sas(t, nameToId)
}*/

func TestSimple(t *testing.T) {
	x := TestStruct{
		Bool: true,
	}
	b := ml_gosn.Marshal(x)

	y := TestStruct{}
	ml_gosn.Unmarshall(b, &y)
	ml_console.PrettyPrint(y)
}

func TestNonStruct(t *testing.T) {
	b := ml_gosn.Marshal("Fuck")
	ml_console.PrintBytes(b, 10)

	m := ""
	ml_gosn.Unmarshall(b, &m)
	fmt.Printf("%v\n", m)
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

func TestTypes(t *testing.T) {
	a := ml_gosn.TypeToString(ml_gosn.T_BOOL)
	if a != "bool" {
		t.Fatal("fuck")
	}
}

func TestGen_forType(t *testing.T) {
	tpl := `
func Test_Type_{{ .Type }}(t *testing.T) {
	a := ml_gosn.TypeToString(ml_gosn.{{ .Type }})
	if a != {{ .TypeString }} {
		t.Fatalf("fuck %v %v", ml_gosn.{{ .Type }}, {{ .TypeString }})
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
	code.Write([]byte("package test_test\n"))
	code.Write([]byte("import (\n\tml_gosn \"github.com/maldan/go-ml/util/encode/gosn\"\n\t\"testing\"\n)"))

	typeList := []string{
		"T_BOOL", "T_8", "T_16", "T_32", "T_64", "T_F32", "T_F64",
		"T_STRING", "T_SHORT_STRING", "T_BIG_STRING",
		"T_SLICE", "T_SHORT_SLICE", "T_BIG_SLICE",
		"T_MAP", "T_SHORT_MAP", "T_BIG_MAP",
		"T_STRUCT", "T_SHORT_STRUCT", "T_BIG_STRUCT",
		"T_CUSTOM", "T_SHORT_CUSTOM", "T_BIG_CUSTOM",
	}
	typeStringList := []any{
		"\"bool\"", "\"uint8\"", "\"uint16\"", "\"uint32\"", "\"uint64\"", "\"float32\"", "\"float64\"",
		"\"string\"", "\"string\"", "\"string\"",
		"\"slice\"", "\"slice\"", "\"slice\"",
		"\"map\"", "\"map\"", "\"map\"",
		"\"struct\"", "\"struct\"", "\"struct\"",
		"\"any\"", "\"any\"", "\"any\"",
	}

	// Type list
	for i, tp := range typeList {
		err = tt.Execute(&code, map[string]any{
			"Type":       tp,
			"TypeString": typeStringList[i],
		})
		if err != nil {
			panic(err)
		}
	}

	ml_file.New("test/types_test.go").Write(code.Bytes())
}

func TestGen_forType2(t *testing.T) {
	tpl := `
func Test_Type2_{{ .FnName }}(t *testing.T) {
	a := ml_gosn.TypeStringToTypeByte({{ .TypeString }})
	if a != ml_gosn.{{ .Type }} {
		t.Fatalf("fuck %v %v", ml_gosn.{{ .Type }}, {{ .TypeString }})
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
	code.Write([]byte("package test_test\n"))
	code.Write([]byte("import (\n\tml_gosn \"github.com/maldan/go-ml/util/encode/gosn\"\n\t\"testing\"\n)"))

	typeStringList := []string{
		"\"bool\"",
		"\"int8\"", "\"uint8\"",
		"\"int16\"", "\"uint16\"",
		"\"int\"", "\"int32\"", "\"uint32\"",
		"\"int64\"", "\"uint64\"",

		"\"float32\"", "\"float64\"",

		"\"string\"",
		"\"slice\"",
		"\"map\"",
		"\"struct\"",
	}
	typeList := []string{
		"T_BOOL",
		"T_8", "T_8",
		"T_16", "T_16",
		"T_32", "T_32", "T_32",
		"T_64", "T_64",
		"T_F32", "T_F64",
		"T_STRING",
		"T_SLICE",
		"T_MAP",
		"T_STRUCT",
	}

	// Type list
	for i, _ := range typeStringList {
		err = tt.Execute(&code, map[string]any{
			"FnName":     strings.ReplaceAll(typeStringList[i], "\"", ""),
			"Type":       typeList[i],
			"TypeString": typeStringList[i],
		})
		if err != nil {
			panic(err)
		}
	}

	ml_file.New("test/types2_test.go").Write(code.Bytes())
}

func TestGen_packAndUnpackSimpleStruct(t *testing.T) {
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
	code.Write([]byte("package test_test\n"))
	code.Write([]byte("import (\n\tml_gosn \"github.com/maldan/go-ml/util/encode/gosn\"\n\t\"testing\"\n)"))

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

	ml_file.New("test/pack_and_unpack_simple_test.go").Write(code.Bytes())
}
