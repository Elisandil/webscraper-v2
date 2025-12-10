package datetime

import (
	"fmt"
	"time"
)

var formats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	time.DateTime,
}

func Parse(dateStr string) (time.Time, error) {

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateStr)
}

func ParseNullable(dateStr string) (*time.Time, error) {

	if dateStr == "" {
		return nil, nil
	}
	t, err := Parse(dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
