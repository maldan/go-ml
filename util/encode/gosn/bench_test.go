package ml_gosn_test

import (
	"encoding/json"
	ml_gosn "github.com/maldan/go-ml/util/encode/gosn"
	"testing"
)

func Benchmark_MarshalGOSN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ml_gosn.Marshal(TestStruct{})
	}
}

func Benchmark_MarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(TestStruct{})
	}
}

func Benchmark_UnmarshalGOSN(b *testing.B) {
	bb := ml_gosn.Marshal(TestStruct{})
	out := TestStruct{}
	for i := 0; i < b.N; i++ {
		ml_gosn.Unmarshall(bb, &out)
	}
}

func Benchmark_UnmarshalJSON(b *testing.B) {
	bb, _ := json.Marshal(TestStruct{})
	out := TestStruct{}
	for i := 0; i < b.N; i++ {
		json.Unmarshal(bb, &out)
	}
}
