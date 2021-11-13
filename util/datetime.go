package util

import "time"

func GetDateTimeString() string {
	return time.Now().Format(time.RFC3339)
}
