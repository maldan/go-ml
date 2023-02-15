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

func TestUpdate(t *testing.T) {
	name := ml_crypto.UID(12)

	userDb := mdb_goson.New[TestUser]("../../../trash/db", name)

	// Insert
	for i := 0; i < 1; i++ {
		userDb.Insert(TestUser{Id: int(userDb.GenerateId()), Username: name, Password: "oglox"})
	}

	// Find
	list := userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Id",
		Limit:     1,
		Where: func(user *TestUser) bool {
			return user.Id == 1
		},
	})
	if !list.IsFound {
		t.Fatalf("fuck")
	}

	// Update
	userDb.UpdateBy(mdb_goson.ArgsUpdate[TestUser]{
		FieldList: "Id",
		Limit:     1,
		Where: func(user *TestUser) bool {
			return user.Id == 1
		},
		Change: func(user *TestUser) {
			user.SalonLogo = "gavno"
		},
	})

	// Find again
	list = userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Id",
		Limit:     1,
		Where: func(user *TestUser) bool {
			return user.Id == 1
		},
	})
	if !list.IsFound {
		t.Fatalf("fuck")
	}
	if list.Result[0].SalonLogo != "gavno" {
		t.Fatalf("update not working")
	}

	// Close
	err := userDb.Close()
	if err != nil {
		panic(err)
	}

	// Delete
	err = ml_file.New("../../../trash/db/" + name).Delete()
	if err != nil {
		panic(err)
	}
}
