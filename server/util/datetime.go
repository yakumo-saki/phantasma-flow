package util

import (
	"time"

	"github.com/yakumo-saki/phantasma-flow/global"
)

var TIME_DEBUG = false

const DEBUG_DATETIME_RFC3339 = "2021-12-25T15:30:45Z09:00"

// @return RFC3339 time string like 2021-12-25T15:30:45Z09:00
func GetDateTimeString() string {
	return GetDateTime().Format(time.RFC3339)
}

// @return yyyymmddhhmmss 20211225153045
func GetDateTimeYyyymmddhhmmss() string {
	return GetDateTime().Format(global.DATETIME_FORMAT)
}

// Do not use time.Now() to unit testability. Use this instead.
// @return time.Now()
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
