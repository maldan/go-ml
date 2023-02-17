package ml_convert_test

import (
	"fmt"
	ml_convert "github.com/maldan/go-ml/util/convert"
	"strconv"
	"testing"
	"time"
)

func TestDataUrl(t *testing.T) {
	data, tp, e := ml_convert.DataUrlToBytes("data:text/plain;base64,SGVsbG8sIFdvcmxkIQ==")
	fmt.Printf("%v\n", string(data))
	fmt.Printf("%v\n", tp)
	fmt.Printf("%v\n", e)
}

func TestE(t *testing.T) {
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		time.Parse("2006-01-01T15:04:05-07:00", "0001-01-01T00:00:00+00:00")
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestToString(t *testing.T) {
	if ml_convert.ToString(1) != "1" {
		t.Fatalf("fuck")
	}
	if ml_convert.ToString(uint8(1)) != "1" {
		t.Fatalf("fuck")
	}
}

func BenchmarkToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if ml_convert.ToString(1) != "1" {
			b.Fatalf("fuck")
		}
	}
}

func BenchmarkStrConv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if strconv.Itoa(1) != "1" {
			b.Fatalf("fuck")
		}
	}
}
