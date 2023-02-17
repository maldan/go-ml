package ml_file_test

import (
	"encoding/json"
	"fmt"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"testing"
)

func TestA(t *testing.T) {
	type tmpStruct struct {
		File ml_file.File
	}
	tmp := tmpStruct{File: *ml_file.New("main_test.go")}

	// Marshal
	f, err := json.Marshal(tmp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", string(f))

	// Unmarshal
	tmp2 := tmpStruct{}
	err = json.Unmarshal(f, &tmp2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", tmp2)
}
