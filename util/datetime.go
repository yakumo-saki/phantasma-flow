package util

import "time"

// @return RFC3339 time string
func GetDateTimeString() string {
	return time.Now().Format(time.RFC3339)
}
