package test_test

import (
	"fmt"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_console "github.com/maldan/go-ml/io/console"
	"testing"
)

func TestCorruptedFind(t *testing.T) {
	name := "corrupted"

	userDb := mdb_goson.New[TestUser]("../../../trash/db", name)

	// userDb.Insert(TestUser{Id: int(userDb.GenerateId()), Username: name, Password: "oglox"})

	// Find
	list := userDb.FindBy(mdb_goson.ArgsFind[TestUser]{
		FieldList: "Id",
		Where: func(user *TestUser) bool {
			return true
		},
	})
	if !list.IsFound {
		t.Fatalf("fuck")
	}

	fmt.Printf("%v\n", list.Count)
}

func TestCorruptedPackAndUnpack(t *testing.T) {
	name := "differedTypes"

	userDb := mdb_goson.New[StructString]("../../../trash/db", name)

	// userDb.Insert(StructNumber{Balance: 10})

	// Find
	list := userDb.FindBy(mdb_goson.ArgsFind[StructString]{
		FieldList: "Balance",
		Where: func(user *StructString) bool {
			return true
		},
	})

	ml_console.PrettyPrint(list.Result)
}
