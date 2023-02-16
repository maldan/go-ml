package ml_time

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"
)

type DateTime struct {
	year  uint16
	month uint8
	day   uint8

	hour   uint8
	minute uint8
	second uint8

	nanoSecond uint16

	tzHour   int8
	tzMinute uint8
}

/*var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}*/
var daysInMonth = [...]uint8{
	31,
	28,
	31,
	30,
	31,
	30,
	31,
	31,
	30,
	31,
	30,
	31,
}

func Now() DateTime {
	return FromTime(time.Now())
}

func FromTime(t time.Time) DateTime {
	// Put location
	loc := t.Format("-07:00")
	lh, _ := strconv.Atoi(loc[1:3])
	lm, _ := strconv.Atoi(loc[4:6])
	if loc[0] == '-' {
		lh = -lh
	}

	return DateTime{
		year:  uint16(t.Year()),
		month: uint8(t.Month()),
		day:   uint8(t.Day()),

		hour:   uint8(t.Hour()),
		minute: uint8(t.Minute()),
		second: uint8(t.Second()),

		nanoSecond: uint16(t.Nanosecond() / 100_000),

		tzHour:   int8(lh),
		tzMinute: uint8(lm),
	}
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	state := "year"
	dotPos := 0
	sign := "+"

	for i := 0; i < len(b); i++ {
		// Date
		if state == "year" && (b[i] == '-' || i == len(b)-1) {
			year, _ := strconv.Atoi(string(b[i-4 : i]))
			d.year = uint16(year)
			state = "month"
			continue
		}
		if state == "month" && (b[i] == '-' || i == len(b)-1) {
			month, _ := strconv.Atoi(string(b[i-2 : i]))
			d.month = uint8(month)
			state = "day"
			continue
		}
		if state == "day" && (b[i] == 'T' || b[i] == ' ' || i == len(b)-1) {
			day, _ := strconv.Atoi(string(b[i-2 : i]))
			d.day = uint8(day)
			state = "hour"
			continue
		}

		// Time
		if state == "hour" && (b[i] == ':' || i == len(b)-1) {
			hour, _ := strconv.Atoi(string(b[i-2 : i]))
			d.hour = uint8(hour)
			state = "minute"
			continue
		}
		if state == "minute" && (b[i] == ':' || i == len(b)-1) {
			minute, _ := strconv.Atoi(string(b[i-2 : i]))
			d.minute = uint8(minute)
			state = "second"
			continue
		}
		if state == "second" && (b[i] == '.' || i == len(b)-1) {
			second, _ := strconv.Atoi(string(b[i-2 : i]))
			d.second = uint8(second)
			state = "nanosecond"
			dotPos = i
			continue
		}
		if state == "nanosecond" && (b[i] == 'Z' || b[i] == ' ' || i == len(b)-1) {
			nsec, _ := strconv.Atoi(string(b[dotPos+1 : i]))
			d.nanoSecond = uint16(nsec)
			state = "timezone"
			continue
		}

		// TimeZone
		if state == "timezone" && (b[i] == '+' || b[i] == '-') {
			sign = string(b[i])
			state = "tzHour"
			continue
		}
		if state == "tzHour" && b[i] == ':' {
			tzHour, _ := strconv.Atoi(string(b[i-2 : i]))
			d.tzHour = int8(tzHour)
			if sign == "-" {
				d.tzHour = -d.tzHour
			}
			state = "tzMinute"
			continue
		}
		if state == "tzMinute" && i == len(b)-1 {
			tzMinute, _ := strconv.Atoi(string(b[i-2 : i]))
			d.tzMinute = uint8(tzMinute)
			state = "end"
			continue
		}
	}

	return nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	sign := "+"
	if d.tzHour < 0 {
		sign = "-"
	}
	return []byte(fmt.Sprintf(
		"\"%04d-%02d-%02dT%02d:%02d:%02d.%dZ%s%02d:%02d\"",
		d.year, d.month, d.day,
		d.hour, d.minute, d.second, d.nanoSecond,
		sign, int(math.Abs(float64(d.tzHour))), d.tzMinute,
	)), nil
}

func (d DateTime) ToBytes() []byte {
	// Prepare
	out := make([]byte, 0, 11)

	// Put date
	out = append(out, uint8(d.year), uint8(d.year>>8))
	out = append(out, d.month, d.day)

	// Put time
	out = append(out, d.hour, d.minute, d.second)
	out = append(out, uint8(d.nanoSecond), uint8(d.nanoSecond>>8))

	// Put time zone
	out = append(out, byte(d.tzHour), d.tzMinute)

	return out
}

func (d *DateTime) FromBytes(b []byte) error {
	// Read date
	d.year = binary.LittleEndian.Uint16(b)
	d.month = b[2]
	d.day = b[3]

	// Read time
	d.hour = b[4]
	d.minute = b[5]
	d.second = b[6]
	d.nanoSecond = binary.LittleEndian.Uint16(b[7:])

	// timeZone
	d.tzHour = int8(b[9])
	d.tzMinute = b[10]

	return nil
}

func (d DateTime) String() string {
	sign := "+"
	if d.tzHour < 0 {
		sign = "-"
	}
	return fmt.Sprintf(
		"%04d-%02d-%02d %02d:%02d:%02d.%05d %s%02d:%02d",
		d.year, d.month, d.day,
		d.hour, d.minute, d.second, d.nanoSecond,
		sign, int(math.Abs(float64(d.tzHour))), d.tzMinute,
	)
}

func (d *DateTime) Year() uint16 {
	return d.year
}
func (d *DateTime) Month() uint8 {
	return d.month
}
func (d *DateTime) Day() uint8 {
	return d.day
}

func (d *DateTime) Hour() uint8 {
	return d.hour
}
func (d *DateTime) Minute() uint8 {
	return d.minute
}
func (d *DateTime) Second() uint8 {
	return d.second
}
func (d *DateTime) Nanosecond() uint16 {
	return d.nanoSecond
}

func (d *DateTime) TimezoneOffset() int {
	return -(int(d.tzHour)*60 + int(d.tzMinute))
}

func (d *DateTime) AddSecond(v int) DateTime {
	nd := *d

	currSec := int(nd.second)
	currMin := int(nd.minute)
	currHour := int(nd.hour)

	daysToOffset := (v + currHour*3600 + currMin*60 + currSec) / 86400

	for i := 0; i < daysToOffset; i++ {
		nd.day += 1

		dim := daysInMonth[nd.month-1]

		// Is leap and february
		if nd.year%4 == 0 && nd.month == 2 {
			// Set 29 days
			dim += 1
		}
		if nd.day > dim {
			nd.day = 1
			nd.month += 1
			if nd.month > 12 {
				nd.month = 1
				nd.year += 1
			}
		}
	}

	// Change time
	nd.hour = uint8((int(nd.hour) + (v+currMin*60+currSec)/3600) % 24)
	nd.minute = uint8((int(nd.minute) + (v+currSec)/60) % 60)
	nd.second = uint8((int(nd.second) + v) % 60)

	return nd
}

func (d *DateTime) In(timeZoneOffset int) DateTime {
	return DateTime{}
}

func (d *DateTime) EqualDate(date *DateTime) bool {
	// Same timezone
	if d.TimezoneOffset() == date.TimezoneOffset() {
		if d.year != date.year {
			return false
		}
		if d.month != date.month {
			return false
		}
		if d.day != date.day {
			return false
		}
	} else {
		// @TODO realize
		return false
	}

	return true
}
