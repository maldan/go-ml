package ml_time

import (
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

func (t Time) String() string {
	return fmt.Sprintf("%v", time.Time(t))
}
