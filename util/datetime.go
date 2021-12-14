package util

import "time"

var TIME_DEBUG = false

const DEBUG_DATETIME_RFC3339 = "2021-12-25T15:30:45Z09:00"

// @return RFC3339 time string
func GetDateTimeString() string {
	if TIME_DEBUG {
		return DEBUG_DATETIME_RFC3339
	}
	return time.Now().Format(time.RFC3339)
}

// @return RFC3339 time string
func GetDateTime() time.Time {
	if TIME_DEBUG {
		t, err := time.Parse(time.RFC3339, DEBUG_DATETIME_RFC3339)
		if err != nil {
			panic(err)
		}
		return t

	}
	return time.Now()
}

// this is pure Now time.
func GetBenchTime() time.Time {
	return time.Now()
}
