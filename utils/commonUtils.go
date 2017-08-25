package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
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

// ParseInt64 to parse int
func ParseInt64(val interface{}) int64 {
	return int64(val.(float64))
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

// GetDigestStr get digest value by string
func GetDigestStr(data string) string {
	return GetDigest(StringToBytes(data))
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

// JSON2String convert json to string
func JSON2String(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

// TrimDir trim directory name by replace "//"
func TrimDir(path string) string {
	return strings.Replace(path, "//", "/", -1)
}

// GetFileSize to get file length
func GetFileSize(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return stat.Size()
}

const (
	//QemuImgTool for virtual size fetching
	QemuImgTool = "/usr/bin/qemu-img"
)

// GetVirtualSize to get file's virtual size
func GetVirtualSize(filepath string) int64 {

	if IsFileExist(QemuImgTool) {
		cmdStr := fmt.Sprintf("%s info %s | grep \"virtual size\" | awk -F' ' '{print $4}' | cut -b2-", QemuImgTool, filepath)
		cmd := exec.Command(cmdStr)
		if err := cmd.Run(); err != nil {
			fmt.Printf("exec cmd %s error\n", cmdStr)
			return 0
		}

		data, err := cmd.Output()
		if err != nil {
			fmt.Printf("get output error\n")
			return 0
		}

		return ParseInt64(string(data))
	}
	return 0
}
