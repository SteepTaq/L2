package date

import (
	"time"
)

func TimeFromString(str string) (time.Time, error) {
	return time.Parse("2006-01-02", str)
}

func StringFromTime(date time.Time) string {
	return date.Format("2006-01-02")
}
