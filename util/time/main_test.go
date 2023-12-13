package ml_time_test

import (
	"encoding/json"
	"fmt"
	ml_time "github.com/maldan/go-ml/util/time"
	"testing"
	"time"
)

var locationList = []string{
	"Europe/Moscow", "UTC", "America/Chicago", "Asia/Calcutta", "Asia/Kathmandu", "Australia/Yancowinna",
	"Canada/Mountain", "Pacific/Marquesas",
}

func TestX(t *testing.T) {
	pp, _ := time.Parse("2006-01-02", "2023-10-11")
	x := time.Since(pp).Hours()
	fmt.Printf("%v\n", x)
}

func TestDateTime(t *testing.T) {
	// Local check
	tm := time.Now()
	today := ml_time.FromTime(tm)
	fmt.Printf("%v\n", today)
	fmt.Printf("%v\n", tm)

	// UTC Check
	tm = time.Now().UTC()
	today = ml_time.FromTime(tm)
	fmt.Printf("%v\n", today)
	fmt.Printf("%v\n", tm)

	// Other
	loc, _ := time.LoadLocation("America/Chicago")
	tm = time.Now().In(loc)
	today = ml_time.FromTime(tm)
	fmt.Printf("%v\n", today)
	fmt.Printf("%v\n", tm)

	// Other
	loc, _ = time.LoadLocation("Asia/Calcutta")
	tm = time.Now().In(loc)
	today = ml_time.FromTime(tm)
	fmt.Printf("%v\n", today)
	fmt.Printf("%v\n", tm)

	// Other
	loc, _ = time.LoadLocation("Asia/Kathmandu")
	tm = time.Now().In(loc)
	today = ml_time.FromTime(tm)
	fmt.Printf("%v\n", today)
	fmt.Printf("%v\n", tm)
}

func TestDateTimeJson(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Kathmandu")

	tm := time.Now().In(loc)
	today := ml_time.FromTime(tm)
	type tmp struct {
		D1 time.Time
		D2 ml_time.DateTime
	}
	testStruct := tmp{
		D1: tm,
		D2: today,
	}

	// Pack
	bytes, _ := json.Marshal(testStruct)
	fmt.Printf("%+v\n", string(bytes))

	// Unpack
	testStruct2 := tmp{}
	json.Unmarshal(bytes, &testStruct2)

	fmt.Printf("%+v\n", testStruct2)
}

func TestDateTimeBytes(t *testing.T) {
	for _, l := range locationList {
		loc, _ := time.LoadLocation(l)

		// Local check
		today := ml_time.FromTime(time.Now().In(loc))
		bytes, _ := today.ToBytes()

		// Unpack
		todayUnpacked := ml_time.DateTime{}
		todayUnpacked.FromBytes(bytes)

		if today.Year() != todayUnpacked.Year() {
			t.Fatalf("fuck")
		}
		if today.Month() != todayUnpacked.Month() {
			t.Fatalf("fuck")
		}
		if today.Day() != todayUnpacked.Day() {
			t.Fatalf("fuck")
		}
		if today.Hour() != todayUnpacked.Hour() {
			t.Fatalf("fuck")
		}
		if today.Minute() != todayUnpacked.Minute() {
			t.Fatalf("fuck")
		}
		if today.Second() != todayUnpacked.Second() {
			t.Fatalf("fuck")
		}
		if today.Nanosecond() != todayUnpacked.Nanosecond() {
			t.Fatalf("fuck")
		}

		if today.TimezoneOffset() != todayUnpacked.TimezoneOffset() {
			t.Fatalf("fuck")
		}

		fmt.Printf("%v - %v - %v\n", today, todayUnpacked, today.TimezoneOffset())
	}
}

func TestDateTimeParseJson(t *testing.T) {
	type tmp struct {
		D2 ml_time.DateTime
	}
	testStruct := tmp{}
	js1 := "{\"D2\":\"2006-01-02 15:04:05.123 +04:12\"}"
	json.Unmarshal([]byte(js1), &testStruct)
	fmt.Printf("%v\n", testStruct.D2)

	testStruct = tmp{}
	js1 = "{\"D2\":\"2006-01-02 15:04:05 +05:12\"}"
	json.Unmarshal([]byte(js1), &testStruct)
	fmt.Printf("%v\n", testStruct.D2)
}

func TestDateTimeFromString(t *testing.T) {
	d1 := ml_time.FromString("2023")
	if d1.Year() != 2023 {
		t.Fatalf("Fuck year %v", d1.Year())
	}
	d1 = ml_time.FromString("2023-03")
	if d1.Month() != 03 {
		t.Fatalf("Fuck month %v", d1.Month())
	}
	d1 = ml_time.FromString("2023-03-05")
	if d1.Day() != 05 {
		t.Fatalf("Fuck day %v", d1.Day())
	}

	d1 = ml_time.FromString("2023-03-05 17:21:10")
	if d1.Second() != 10 {
		t.Fatalf("Fuck second %v", d1.Second())
	}

	d1 = ml_time.FromString("2023-03-05 17:21:10.123")
	if d1.Nanosecond() != 123 {
		t.Fatalf("Fuck nano second %v", d1.String())
	}

	d1 = ml_time.FromString("2023-03-05 17:21:10.123 +01:00")
	if d1.TimezoneOffset() != -60 {
		t.Fatalf("Fuck tz %v", d1.TimezoneOffset())
	}

	d1 = ml_time.FromString("2023-03-05 17:21:10 +01:00")
	if d1.TimezoneOffset() != -60 {
		t.Fatalf("Fuck tz %v", d1.TimezoneOffset())
	}

	// Hz
	d1 = ml_time.FromString("2023-12-12T23:42:00Z")
	if d1.Year() != 2023 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.Month() != 12 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.Day() != 12 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.Hour() != 23 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.Minute() != 42 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.Second() != 0 {
		t.Fatalf("Fuck %v", d1.String())
	}
	if d1.TimezoneOffset() != 0 {
		t.Fatalf("Fuck %v", d1.String())
	}
	fmt.Printf("%v\n", d1.String())
}

func TestDateTimeAdd(t *testing.T) {
	t1 := time.Now()
	t2 := ml_time.FromTime(t1)

	gapList := []int{1, 100, 1000}

	for _, gap := range gapList {
		timing := time.Now()
		for i := -1_000_000; i < 1_000_000; i++ {
			tn1 := t1.Add(time.Second * time.Duration(i*gap))
			tn2 := t2.AddSecond(i * gap)

			if tn1.Second() != int(tn2.Second()) {
				t.Fatalf("second fuck %v %v", tn1, tn2)
			}
			if tn1.Minute() != int(tn2.Minute()) {
				t.Fatalf("minute fuck %v %v", tn1, tn2)
			}
			if tn1.Hour() != int(tn2.Hour()) {
				t.Fatalf("hour fuck %v %v", tn1, tn2)
			}

			if tn1.Day() != int(tn2.Day()) {
				t.Fatalf("day fuck\n%v\n%v\n%v", tn1, tn2, i)
			}
			if int(tn1.Month()) != int(tn2.Month()) {
				t.Fatalf("month fuck")
			}
			if tn1.Year() != int(tn2.Year()) {
				t.Fatalf("year fuck")
			}
		}
		fmt.Printf("%v\n", time.Since(timing))
	}
}

func TestDateTimeAdd2(t *testing.T) {
	t1 := time.Now()
	t2 := ml_time.FromTime(t1)

	fmt.Printf("%v\n", t1.Add(time.Second*86400*4200))
	fmt.Printf("%v\n", t2.AddSecond(86400*4200))

	fmt.Printf("NOW: %v\n", t1.Format("2006-01-02 15:04:05"))
	fmt.Printf("NED: %v\n", t1.Add(time.Second*time.Duration(1_000_000*5000)).Format("2006-01-02 15:04:05"))
	fmt.Printf("HAV: %v\n", t2.AddSecond(1_000_000*5000))

	fmt.Printf("NOW: %v\n", t1.Format("2006-01-02 15:04:05"))
	fmt.Printf("NED: %v\n", t1.Add(time.Second*time.Duration(-1_000_000*5000)).Format("2006-01-02 15:04:05"))
	fmt.Printf("HAV: %v\n", t2.AddSecond(-1_000_000*5000))
}

func TestDateTimeEqual(t *testing.T) {
	t1 := time.Now()
	t2 := ml_time.FromTime(t1)

	fmt.Printf("%v\n", t2.In(0))
	fmt.Printf("%v\n", t2.UTC())
	fmt.Printf("%v\n", t2.UTC().In(-180))

	fmt.Printf("%v\n", t2.UTC().EqualDate(t2))
	fmt.Printf("%v\n", t2.UTC().EqualTime(t2))
}

func TestDateTimeToTime(t *testing.T) {
	t1 := time.Now()
	t2 := ml_time.FromTime(t1)
	t3 := t2.ToTime()

	fmt.Printf("%v\n", t1)
	fmt.Printf("%v\n", t2)
	fmt.Printf("%v\n", t3)
}
