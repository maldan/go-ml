package test_test

import (
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"

	ml_console "github.com/maldan/go-ml/util/io/console"
	"testing"
)

type TestStruct struct {
	Name string
}

func TestInsert(t *testing.T) {
	db := mdb.New[TestStruct]("../../../trash/db", "gas", &gosn_driver.Container{})
	db.Insert(TestStruct{Name: "A"})
}

func TestFind(t *testing.T) {
	db := mdb.New[TestStruct]("../../../trash/db", "gas", &gosn_driver.Container{})
	db.Insert(TestStruct{Name: "A"})

	// Find
	list := db.FindBy(mdb.ArgsFind[TestStruct]{
		FieldList: "Name",
		Where: func(user *TestStruct) bool {
			return true
		},
	})

	ml_console.PrettyPrint(list.Result)
}
