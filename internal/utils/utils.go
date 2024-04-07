package utils

import "time"

func CurrentTimeString() string {
	return time.Now().Format(time.RFC3339Nano)
}
