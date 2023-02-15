package test_test

import (
	"fmt"
	ml_crypto "github.com/maldan/go-ml/crypto"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_file "github.com/maldan/go-ml/io/fs/file"
	"testing"
	"time"
)

func TestInsertMeasureTime_1000(t *testing.T) {
	name := ml_crypto.UID(12)

	userDb := mdb_goson.New[TestUser]("../../../trash/db", name)

	// Insert
	tt := time.Now()
	for i := 0; i < 1000; i++ {
		userDb.Insert(TestUser{Id: int(userDb.GenerateId()), Username: name, Password: "oglox"})
	}
	fmt.Printf("%v\n", time.Since(tt))

	err := userDb.Close()
	if err != nil {
		panic(err)
	}

	err = ml_file.New("../../../trash/db/" + name).Delete()
	if err != nil {
		panic(err)
	}
}
