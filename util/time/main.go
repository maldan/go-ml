package ml_time

import (
	"fmt"
	"strings"
	"time"
)

type Time time.Time

var timeParseTemplateList = []string{
	"2006-01-02 15:04:05", "2006-01-02", "2006-01-02 15:04", "2006-01-02 15:04:05.999",
}

func (t *Time) UnmarshalJSON(b []byte) error {
	// Get rid of "
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

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

func (t Time) String() string {
	return fmt.Sprintf("%v", time.Time(t))
}
