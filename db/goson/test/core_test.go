package test_test

import (
	"github.com/maldan/go-ml/db/goson/core"
	ml_slice "github.com/maldan/go-ml/slice"
	"testing"
)

func TestGetNameList(t *testing.T) {
	nl := core.GetNameList(TestA{})
	// fmt.Printf("%v\n", nl)
	if len(nl) != 2 {
		t.Fatalf("Fuck!")
	}
	if !ml_slice.Includes(nl, "Name") {
		t.Fatalf("Fuck!")
	}
	if !ml_slice.Includes(nl, "Crazy") {
		t.Fatalf("Fuck!")
	}
}
