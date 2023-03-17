package test_test

import (
	"fmt"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"
	"time"

	ml_console "github.com/maldan/go-ml/util/io/console"
	"testing"
)

type TestStruct struct {
	Name string
}

func TestInsert(t *testing.T) {
	db := mdb.New[TestStruct]("../../../trash/db", "gas", &gosn_driver.Container{})
	for i := 0; i < 10000; i++ {
		db.Insert(TestStruct{Name: "A"})
	}
}

func TestFind(t *testing.T) {
	db := mdb.New("../../../trash/db", "gas", TestStruct{}, &gosn_driver.Container{})
	db.Insert(TestStruct{Name: "A"})

	// Find
	tt := time.Now()
	list := db.FindBy(mdb.ArgsFind{
		FieldList: "Name",
		Where: func(user any) bool {
			return user.(*TestStruct).Name == "A"
		},
	})
	fmt.Printf("%v\n", time.Since(tt))

	ml_console.PrettyPrint(list.Result)
}

func TestFind2(t *testing.T) {
	db := mdb.New[TestStruct]("../../../trash/db", "gas", &gosn_driver.Container{})
	db.Insert(TestStruct{Name: "A"})

	// Find
	tt := time.Now()
	list := db.FindBy(mdb.ArgsFind[TestStruct]{
		FieldList: "Name",
		Where: func(user *TestStruct) bool {
			return user.Name == "A"
		},
	})
	fmt.Printf("%v\n", time.Since(tt))

	ml_console.PrettyPrint(list.Result)
}
