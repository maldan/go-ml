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
var daysInMonth = [...]int{
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

func (d *DateTime) SetTimezoneOffset(timeZoneOffset int) {
	d.tzHour = -int8(timeZoneOffset/60) % 60
	d.tzMinute = uint8(timeZoneOffset % 60)
}

func (d *DateTime) AddSecond(v int) DateTime {
	nd := *d

	// Offset forward of backward
	dayOffset := 1
	if v < 0 {
		dayOffset = -1
	}

	// Total days to offset
	totalSec := int(nd.hour)*3600 + int(nd.minute)*60 + int(nd.second)
	daysToOffset := (totalSec + v) / 86400
	if daysToOffset < 0 {
		daysToOffset = -(daysToOffset)
	}

	// If now today time is 100 for example, and I subtract 101. Then time will be -1.
	// But -1 means 86399. So it means we need go back to 1 day.
	dede := (totalSec + v) % 86400
	if dede < 0 {
		daysToOffset += 1
	}

	// New day
	newDay := int(nd.day)
	newMonth := int(nd.month)
	newYear := int(nd.year)

	if dayOffset > 0 {
		for i := 0; i < daysToOffset; i++ {
			newDay += 1

			dim := daysInMonth[newMonth-1]

			// Is leap and february
			if ((newYear%4 == 0 && newYear%100 != 0) || newYear%400 == 0) && newMonth == 2 {
				// Set 29 days
				dim += 1
			}
			if newDay > dim {
				newDay = 1
				newMonth += 1
				if newMonth > 12 {
					newMonth = 1
					newYear += 1
				}
			}
		}
	} else {
		for i := 0; i < daysToOffset; i++ {
			newDay -= 1

			// Backward
			if newDay < 1 {
				newMonth -= 1
				if newMonth < 1 {
					newMonth = 12
					newYear -= 1
				}

				// Is leap and february
				if ((newYear%4 == 0 && newYear%100 != 0) || newYear%400 == 0) && newMonth == 2 {
					newDay = daysInMonth[newMonth-1] + 1
				} else {
					newDay = daysInMonth[newMonth-1]
				}
			}
		}
	}

	// Set new date
	nd.year = uint16(newYear)
	nd.month = uint8(newMonth)
	nd.day = uint8(newDay)

	// Set new time
	gav := ((int(nd.hour)*3600 + int(nd.minute)*60 + int(nd.second)) + v) % 86400
	if gav < 0 {
		gav += 86400
	}

	nd.hour = uint8((gav / 3600) % 24)
	nd.minute = uint8((gav % 3600) / 60)
	nd.second = uint8(gav % 60)

	return nd
}

// In timeZoneOffset is offset in minutes. For example -180 means +03:00.
func (d DateTime) In(timeZoneOffset int) DateTime {
	// Difference between current timezone and given
	diff := d.TimezoneOffset() - timeZoneOffset
	offset := d.AddSecond(diff * 60)
	offset.SetTimezoneOffset(timeZoneOffset)
	return offset
}

func (d DateTime) UTC() DateTime {
	return d.In(0)
}

// Equal compares date and time. For example
// 2006-01-02 15:04:03 == 1992-01-02 12:33:11
func (d DateTime) Equal(date DateTime) bool {
	return d.EqualDate(date) && d.EqualTime(date)
}

// PreciseEqual same as Equal but also compare nanoseconds
func (d DateTime) PreciseEqual(date DateTime) bool {
	return d.Equal(date) && d.nanoSecond == date.nanoSecond
}

// EqualDate compare only dates. For example
// 2006-01-02 == 1992-01-02
func (d DateTime) EqualDate(date DateTime) bool {
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
		// Make dates equal with timezone
		dateIn := date.In(d.TimezoneOffset())

		if d.year != dateIn.year {
			return false
		}
		if d.month != dateIn.month {
			return false
		}
		if d.day != dateIn.day {
			return false
		}
	}

	return true
}

// EqualTime compare only time. For example
// 15:32:12 == 12:11:33
func (d DateTime) EqualTime(date DateTime) bool {
	// Same timezone
	if d.TimezoneOffset() == date.TimezoneOffset() {
		if d.hour != date.hour {
			return false
		}
		if d.minute != date.minute {
			return false
		}
		if d.second != date.second {
			return false
		}
	} else {
		// Make dates equal with timezone
		dateIn := date.In(d.TimezoneOffset())

		if d.hour != dateIn.hour {
			return false
		}
		if d.minute != dateIn.minute {
			return false
		}
		if d.second != dateIn.second {
			return false
		}
	}

	return true
}
