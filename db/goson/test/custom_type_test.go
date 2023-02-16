package test_test

import (
	"fmt"
	ml_crypto "github.com/maldan/go-ml/crypto"
	mdb_goson "github.com/maldan/go-ml/db/goson"
	ml_file "github.com/maldan/go-ml/io/fs/file"
	ml_time "github.com/maldan/go-ml/util/time"
	"testing"
	"time"
)

func TestCustomTypeFind(t *testing.T) {
	name := ml_crypto.UID(12)

	userDb := mdb_goson.New[StructCustomType]("../../../trash/db", name)

	// Insert
	values := make([]StructCustomType, 0)
	for i := 0; i < 1_000_000; i++ {
		values = append(values, StructCustomType{
			Time: ml_time.Time(time.Now().Add(time.Hour * time.Duration(i))),
		})
	}
	lastDay := ml_time.Time(time.Now().Add(time.Hour * time.Duration(1_000_100)))

	tm := time.Now()
	userDb.InsertMany(values)
	userDb.Insert(StructCustomType{Time: lastDay})
	fmt.Printf("Insert Many Time: %v\n", time.Since(tm))

	// Find
	tm = time.Now()
	list := userDb.FindBy(mdb_goson.ArgsFind[StructCustomType]{
		FieldList: "Time",
		Limit:     1,
		Where: func(user *StructCustomType) bool {
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
}
