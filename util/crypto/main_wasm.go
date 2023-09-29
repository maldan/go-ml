package ml_crypto

import (
	"crypto/rand"
	"math/big"
	"time"
)

var a = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
var a2 = "0123456789"
var al = len(a)
var al2 = len(a2)
var b = big.NewInt(1_000_000_000)

func UID(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix()

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a[int(num.Int64()+t)%al]
	}

	return string(out)
}

func RandomNumberCode(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix()

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a2[int(num.Int64()+t)%al2]
	}

	return string(out)
}
