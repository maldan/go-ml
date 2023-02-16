package ml_time

import (
	"encoding/binary"
	"fmt"
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

	out := make([]byte, 0, 8)
	out = append(out, uint8(year))
	out = append(out, uint8(year>>8))

	out = append(out, uint8(tm.Month()))
	out = append(out, uint8(tm.Day()))

	return out
}

func (t *Time) FromBytes(b []byte) error {
	tm := time.Time{}

	year := binary.LittleEndian.Uint16(b)
	month := b[2]
	day := b[3]

	*t = Time(time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, tm.Location()))

	return nil
}

func (t Time) String() string {
	return fmt.Sprintf("%v", time.Time(t))
}
