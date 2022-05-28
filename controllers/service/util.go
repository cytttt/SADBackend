package service

import "time"

// period defintion:
// morning   06:00 ~ 12:00
// afternoon 12:00 ~ 18:00
// evening   18:00 ~ 21:00
// midnight  21:00 ~ 24:00
type timeRange string

const (
	TIME_MORNING   timeRange = "morning"
	TIME_AFTERNOON timeRange = "afternoon"
	TIME_EVENING   timeRange = "evening"
	TIME_MIDNIGHT  timeRange = "midnight"
)

func string2Time(timeStr, format string) (*time.Time, error) {
	offset := int((8 * time.Hour).Seconds())
	loc := time.FixedZone("Asia/Taipei", offset)
	newTime, err := time.ParseInLocation(format, timeStr, loc)
	if err != nil {
		return nil, err
	}
	return &newTime, err
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
