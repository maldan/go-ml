package ml_hash_test

import (
	ml_hash "github.com/maldan/go-ml/util/hash"
	"testing"
)

func TestSha1(t *testing.T) {
	if ml_hash.Sha1("sex") != "68bb04bd54b8f6c530695e0b77de298276a0511d" {
		t.Errorf("Sha1 not working correctly")
	}
	if ml_hash.Sha1([]byte{97}) != "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8" {
		t.Errorf("Sha1 not working correctly")
	}
}

func TestSha256(t *testing.T) {
	if ml_hash.Sha256("sex") != "98d44e13f455d916674d38424d39e1cb01b2a9132aacbb7b97a6f8bb7feb2544" {
		t.Errorf("Sha256 not working correctly")
	}
	if ml_hash.Sha256([]byte{97}) != "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb" {
		t.Errorf("Sha256 not working correctly")
	}
}
