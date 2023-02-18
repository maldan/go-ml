package test_test

import (
	"fmt"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_crypto "github.com/maldan/go-ml/util/crypto"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"

	ml_time "github.com/maldan/go-ml/util/time"
	"testing"
	"time"
)

/*func TestCustomTypeFind(t *testing.T) {
	// Define type
	type structCustomType struct {
		Time ml_time.Time
	}

	// Create DB
	name := ml_crypto.UID(12)
	userDb := mdb_goson.New[structCustomType]("../../../trash/db", name)

	// Insert
	values := make([]structCustomType, 0)
	for i := 0; i < 1_000_000; i++ {
		values = append(values, structCustomType{
			Time: ml_time.Time(time.Now().Add(time.Hour * time.Duration(i))),
		})
	}
	lastDay := ml_time.Time(time.Now().Add(time.Hour * time.Duration(1_000_100)))

	tm := time.Now()
	userDb.InsertMany(values)
	userDb.Insert(structCustomType{Time: lastDay})
	fmt.Printf("Insert Many Time: %v\n", time.Since(tm))

	// Find
	tm = time.Now()
	list := userDb.FindBy(mdb_goson.ArgsFind[structCustomType]{
		FieldList: "Time",
		Limit:     1,
		Where: func(user *structCustomType) bool {
			d, m, y := time.Time(user.Time).Date()
			d1, m1, y1 := time.Time(lastDay).Date()
			return d == d1 && m == m1 && y == y1
		},
	})
	fmt.Printf("Find Time: %v\n", time.Since(tm))
	if !list.IsFound {
		t.Fatalf("fuck")
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
}*/

func TestCustomTypeFind2(t *testing.T) {
	// Define type
	type structCustomType struct {
		Time ml_time.DateTime
	}

	// Create DB
	name := ml_crypto.UID(12)
	userDb := mdb_goson.New[structCustomType]("../../../trash/db", name, &gosn_driver.Container{})

	// Insert 88 ms vs 39 ms
	values := make([]structCustomType, 0)
	for i := 0; i < 10; i++ {
		values = append(values, structCustomType{
			Time: ml_time.Now().AddSecond(i),
		})
	}
	lastDay := ml_time.Now().AddSecond(1_000_100)

	tm := time.Now()
	userDb.InsertMany(values)
	userDb.Insert(structCustomType{Time: lastDay})
	fmt.Printf("Insert Many Time: %v\n", time.Since(tm))

	// Find
	tm = time.Now()
	list := userDb.FindBy(mdb_goson.ArgsFind[structCustomType]{
		FieldList: "Time",
		Limit:     1,
		Where: func(user *structCustomType) bool {
			//d, m, y := time.Time(user.Time).Date()
			//d1, m1, y1 := time.Time(lastDay).Date()
			//return d == d1 && m == m1 && y == y1
			return lastDay.EqualDate(user.Time)
		},
	})
	fmt.Printf("Find Time: %v\n", time.Since(tm))
	if !list.IsFound {
		t.Fatalf("fuck")
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
