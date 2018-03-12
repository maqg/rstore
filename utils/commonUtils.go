package utils

import (
	"octlink/rstore/utils/configuration"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"octlink/rstore/modules/config"
	"octlink/rstore/utils/octlog"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
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

// Time2StrSimple convert int64 time to string
func Time2StrSimple(timeVal int64) string {
	return time.Unix(timeVal, 0).Format("20060102150405")
}

// CurrentTime get current time in int64 format
func CurrentTime() int64 {
	return int64(time.Now().Unix())
}

// CurrentTimeStr get current time in string format
func CurrentTimeStr() string {
	return Time2Str(CurrentTime())
}

// CurrentTimeSimple get current time in string format
func CurrentTimeSimple() string {
	return Time2StrSimple(CurrentTime())
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

// JSON2Bytes convert json to string
func JSON2Bytes(v interface{}) []byte {
	data, _ := json.MarshalIndent(v, "", "  ")
	return data
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

// GetFileData to get file content by index and length
func GetFileData(filepath string, index int, length int) ([]byte, error) {

	if !IsFileExist(filepath) {
		return nil, errors.New("file not exist")
	}

	fd, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("open file error")
	}

	defer fd.Close()

	fd.Seek(int64(configuration.BlobSize * index), 0)
	buffer := make([]byte, length)
	n, err := fd.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n != length {
		return nil, errors.New("no enough bytes data to fetch")
	}

	return buffer, nil
}

const (
	//QemuImgTool for virtual size fetching
	QemuImgTool = "/usr/bin/qemu-img"
)

// OCTSystem for syscal command calling
func OCTSystem(cmdstr string) (string, error) {

	cmd := exec.Command("/bin/sh", "-c", cmdstr)
	data, err := cmd.Output()
	if err != nil {
		fmt.Printf("get cmd output error of %s,%s\n", cmdstr, err)
		return "", err
	}

	return BytesToString(data), nil
}

// GetVirtualSize to get file's virtual size
func GetVirtualSize(filepath string) int64 {

	if !IsFileExist(QemuImgTool) {
		return 0
	}

	cmdStr := fmt.Sprintf("%s info %s | grep \"virtual size\" | awk -F' ' '{print $4}' | cut -b2-",
		QemuImgTool, filepath)

	result, err := OCTSystem(cmdStr)
	if err != nil {
		fmt.Printf("exec cmd error %s:%s\n", cmdStr, err)
		return 0
	}

	return StringToInt64(result)
}

// StringToInt convert string to int value
func StringToInt(src string) int {
	ret, err := strconv.Atoi(src)
	if err != nil {
		return -1
	}
	return ret
}

// StringToInt64 convert string to int value
func StringToInt64(src string) int64 {
	src = strings.Replace(src, "\n", "", -1)
	ret, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return -1
	}
	return ret
}

// IntToString convert int to string value
func IntToString(src int) string {
	return strconv.Itoa(src)
}

// Int64ToString convert int64 to string value
func Int64ToString(src int64) string {
	return strconv.FormatInt(src, 10)
}

// FileToBytes for filepath convert to bytes
func FileToBytes(filepath string) []byte {
	if !IsFileExist(filepath) {
		return nil
	}

	fd, err := os.Open(filepath)
	if err != nil {
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil
	}

	return data
}

// FileToString convert file content to string
func FileToString(filepath string) string {
	return BytesToString(FileToBytes(filepath))
}

// NumberToInt convert int,int32,int64,float,float32, to int
func NumberToInt(value interface{}) int {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return value.(int)
	case reflect.Int64:
		return int(value.(int64))
	case reflect.Float32:
		return int(value.(float32))
	case reflect.Float64:
		return int(value.(float64))
	}
	return 0
}

// SendUserSignal send USR1 Signal to process name
func SendUserSignal(pname string) {
	if !IsPlatformWindows() {
		cmd := fmt.Sprintf("pidof %s | xargs kill -USR1 > /dev/null 2>&1", pname)
		OCTSystem(cmd)
	}
}

// String2Int convert string to 32-bit int
func String2Int(src string) int {
	val, err := strconv.ParseInt(src, 10, 32)
	if err != nil {
		return -1
	}
	return int(val)
}

// String2Int64 convert string to int64
func String2Int64(src string) int64 {
	val, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return -1
	}
	return val
}

// ParseBlobDigest for common digest, return dgst, 0, 0
// for huge blob, return dgst, index, length
func ParseBlobDigest(dgst string) (string, int, int) {
	segs := strings.Split(dgst, "_")
	if len(segs) == 1 {
		return dgst, 0, 0
	}

	return segs[0], String2Int(segs[1]), String2Int(segs[2])
}

// IsHugeBlobDigest if hugeblobdigest like xxxxx_2_33333, return ture
func IsHugeBlobDigest(dgst string) bool {
	segs := strings.Split(dgst, "_")
	return len(segs) == 1
}

// CopyFile for srcfile to dst file, return size on success
func CopyFile(srcFile, dstFile string) (int64, error) {
	sd, err := os.Open(srcFile)
	if err != nil {
		octlog.Error("open src file of %s error: %s\n", srcFile, err)
		return 0, err
	}

	defer sd.Close()

	dd, err := os.Create(dstFile)
	if err != nil {
		octlog.Error("open dst file of %s error: %s\n", dstFile, err)
		return 0, err
	}
	defer dd.Close()

	return io.Copy(dd, sd)
}

// OSType return os type
func OSType() string {
	// can be darwin,windows,linux
	return runtime.GOOS
}

// IsPlatformWindows for platform type judgement
func IsPlatformWindows() bool {
	return OSType() == config.OSTypeWindows
}

var logger *octlog.LogConfig

// InitLog to init api log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("utils.log", level)
}

// RemoveFromArray for remove from array
func RemoveFromArray(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}

// RemoveFromArrayEx for remove from array without seq indexing
func RemoveFromArrayEx(s []interface{}, i int) []interface{} {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
