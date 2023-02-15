package test_test

import (
	"fmt"
	"testing"
	"time"
)

/*func TestMy1(t *testing.T) {
	s := "сука"
	for i := 0; i < len(s); i++ {
		fmt.Printf("%v\n", string(s[i]))
	}
}*/

/*func TestWriteAndRead(t *testing.T) {
	table := cdb_goson.New[Test]("../db/test")

	rndName := cmhp_crypto.UID(8)

	table.Insert(Test{
		FirstName: fmt.Sprintf(rndName),
	})

	rs := table.FindBy(cdb_goson.ArgsFind[Test]{
		FieldList: "FirstName",
		Where: func(test *Test) bool {
			return test.FirstName == rndName
		},
	})

	if !rs.IsFound {
		t.Fatalf("fuck")
	}
}*/

func TestX3(t *testing.T) {
	a := 0
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		a = 0b0000_1111
		for j := 0; j < 64; j++ {
			a = a << 1
		}
	}
	fmt.Printf("T2: %v\n", time.Since(tt))
	fmt.Printf("%v\n", a)
}

func TestX2(t *testing.T) {
	type Gas struct {
		Id int
		S  []string
	}

	xyz := make([]Gas, 1_000_000)
	for i := 0; i < 1_000_000; i++ {
		xyz[i].Id = 1_000_000 - 1
	}

	a := 0
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		a += xyz[i].Id
	}
	fmt.Printf("T2: %v\n", time.Since(tt))
	fmt.Printf("%v\n", a)
	/*fmt.Printf("%v\n", xyz[0])

	tt := time.Now()
	sort.SliceStable(xyz, func(i, j int) bool {
		return xyz[i].Id > xyz[j].Id
	})
	fmt.Printf("T2: %v\n", time.Since(tt))

	fmt.Printf("%v\n", xyz[0])

	for {
		time.Sleep(time.Second)
	}*/
}
