package ml_time

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Time time.Time

var timeParseTemplateList = []string{
	"2006-01-02T15:04:05.999-07:00",
	"2006-01-02 15:04:05.999-07:00",
	"2006-01-02 15:04:05.999",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006-01",
	"2006",
}

func (t *Time) UnmarshalJSON(b []byte) error {
	// Get rid of "
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	l, _ := time.LoadLocation("Etc/GMT")
	fmt.Printf("%v\n", l)
	fmt.Printf("%v\n", time.Now().In(l))

	// Parse time
	for _, tpl := range timeParseTemplateList {
		pt, err := time.Parse(tpl, value)
		if err == nil {
			// Set result
			*t = Time(pt)
			return nil
		}
	}

	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	tt := time.Time(t)

	out := ""
	out += tt.Format("2006-01-02T15:04:05")
	ms := tt.Format(".999")
	if ms == "" {
		out += ".000"
	} else {
		out += ms
	}
	out += tt.Format("-07:00")

	return []byte("\"" + out + "\""), nil
}

func (t Time) ToBytes() []byte {
	tm := time.Time(t)

	year := tm.Year()

	// Put date
	out := make([]byte, 0, 12)
	out = append(out, uint8(year), uint8(year>>8))
	out = append(out, uint8(tm.Month()), uint8(tm.Day()))

	// Put time
	hour := tm.Hour()
	minute := tm.Minute()
	second := tm.Second()
	nsec := tm.Nanosecond() / 10
	out = append(out, uint8(hour), uint8(minute), uint8(second), uint8(nsec>>8), uint8(nsec))

	// Put location
	loc := tm.Format("-07:00")
	lh, _ := strconv.Atoi(loc[1:3])
	lm, _ := strconv.Atoi(loc[4:6])
	if loc[0] == '-' {
		lh = -lh
	}
	out = append(out, uint8(lh), uint8(lm))

	return out
}

func (t *Time) FromBytes(b []byte) error {
	tm := time.Time{}

	fmt.Printf("%v\n", len(b))

	year := binary.LittleEndian.Uint16(b)
	month := b[2]
	day := b[3]

	*t = Time(time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, tm.Location()))

	return nil
}

func (t Time) String() string {
	return fmt.Sprintf("%v", time.Time(t))
}
