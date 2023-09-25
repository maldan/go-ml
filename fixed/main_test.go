package mfixed_test

import (
	"fmt"
	mfixed "github.com/maldan/go-ml/fixed"
	"testing"
)

func TestS(t *testing.T) {
	a := mfixed.New(10, 3000)
	b := mfixed.New(32, 6000)
	fmt.Printf("%v\n", a.Mul(b).ToString())
	fmt.Printf("%v\n", 10.3*32.6)
}
