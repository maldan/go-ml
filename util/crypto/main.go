package ml_crypto

import (
	"crypto/rand"
	"math/big"
	"os"
	"runtime"
	"time"
)

var m runtime.MemStats
var a = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
var a2 = "0123456789"
var al = len(a)
var al2 = len(a2)
var b = big.NewInt(1_000_000_000)

func init() {
	runtime.ReadMemStats(&m)
}

func UID(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix() + int64(os.Getpid()) + int64(m.TotalAlloc+m.Alloc+m.Sys)

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a[int(num.Int64()+t)%al]
	}

	return string(out)
}

func RandomNumberCode(size int) string {
	out := make([]byte, size)
	t := time.Now().Unix() + int64(os.Getpid()) + int64(m.TotalAlloc+m.Alloc+m.Sys)

	for i := 0; i < size; i++ {
		num, _ := rand.Int(rand.Reader, b)
		out[i] = a2[int(num.Int64()+t)%al2]
	}

	return string(out)
}

/*func EncryptAes32[T string | []byte](data T, password string) (T, error) {
	// Prepare key
	key := make([]byte, 32)
	h := ml_hash.Sha1(password)
	for i := 0; i < len(h); i++ {
		key[i%32] += h[i]
	}

	// Prepare cipher
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return *new(T), err
	}

	// Gcm
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return *new(T), err
	}

	// Read rand
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return *new(T), err
	}

	switch any(data).(type) {
	case string:
		v := gcm.Seal(nonce, nonce, []byte(data), nil)
		return any(cmhp_convert.ToBase64(v)).(T), nil
	default:
		v := gcm.Seal(nonce, nonce, []byte(data), nil)
		return any(v).(T), nil
	}
}

func DecryptAes32[T string | []byte](data T, password string) (T, error) {
	// Prepare key
	key := make([]byte, 32)
	h := ml_hash.Sha1(password)
	for i := 0; i < len(h); i++ {
		key[i%32] += h[i]
	}

	// Prepare cipher
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return *new(T), err
	}

	// Gcm
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return *new(T), err
	}

	switch any(data).(type) {
	case string:
		data2 := cmhp_convert.FromBase64(any(data).(string))
		nonce, ciphertext := data2[:gcm.NonceSize()], data2[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		str := string(plaintext)
		return any(str).(T), err
	default:
		data2 := any(data).([]byte)
		nonce, ciphertext := data2[:gcm.NonceSize()], data2[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		return any(plaintext).(T), err
	}
}*/
