package ml_hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type hashable interface {
	string | []byte
}

func _hasher[T hashable](h hash.Hash, data T) string {
	switch any(data).(type) {
	case []byte:
		h.Write(any(data).([]byte))
		return hex.EncodeToString(h.Sum(nil))
	case string:
		h.Write([]byte(any(data).(string)))
		return hex.EncodeToString(h.Sum(nil))
	default:
		panic("unsupported type")
	}
}

func Sha1[T hashable](data T) string {
	return _hasher(sha1.New(), data)
}

func Sha256[T hashable](data T) string {
	return _hasher(sha256.New(), data)
}
