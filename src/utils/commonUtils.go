package utils

import (
	"os"
	"time"
)

const (
	TIME_STR_FORMAT = "2006-01-02 15:04:05"
)

func Time2Str(timeVal int64) string {
	return time.Unix(timeVal, 0).Format(TIME_STR_FORMAT)
}

func Version() string {
	return "0.0.1"
}

func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ParasInt(val interface{}) int {
	return int(val.(float64))
}
