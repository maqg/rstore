package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"reflect"
	"time"
	"unsafe"
)

const (
	timeStrFormat = "2006-01-02 15:04:05"
)

// Time2Str convert int64 time to string
func Time2Str(timeVal int64) string {
	return time.Unix(timeVal, 0).Format(timeStrFormat)
}

// CurrentTime get current time in int64 format
func CurrentTime() int64 {
	return int64(time.Now().Unix())
}

// CurrentTimeStr get current time in string format
func CurrentTimeStr() string {
	return Time2Str(CurrentTime())
}

// Version of this program
func Version() string {
	return "0.0.1"
}

// IsFileExist check file's existence
func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// ParseInt to parse int
func ParseInt(val interface{}) int {
	return int(val.(float64))
}

// StringToBytes return GoString's buffer slice(enable modify string)
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToString convert b to string without copy
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToPointer returns &s[0], which is not allowed in go
func StringToPointer(s string) unsafe.Pointer {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(p.Data)
}

// BytesToPointer returns &b[0], which is not allowed in go
func BytesToPointer(b []byte) unsafe.Pointer {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(p.Data)
}

// GetDigest return digest hash value of sha256
func GetDigest(buffer []byte) string {
	hash := sha256.New()
	hash.Write(buffer)
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateDir create dir if not exist
func CreateDir(filepath string) {
	if !IsFileExist(filepath) {
		os.MkdirAll(filepath, os.ModePerm)
	}
}

// RemoveDir if file or directory exists, just remove it
func RemoveDir(filepath string) {
	if IsFileExist(filepath) {
		os.RemoveAll(filepath)
	}
}

// Remove if file or directory exists, just remove it
func Remove(filepath string) {
	if IsFileExist(filepath) {
		os.Remove(filepath)
	}
}
