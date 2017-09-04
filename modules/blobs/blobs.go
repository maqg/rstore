package blobs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"os"
)

var logger *octlog.LogConfig

// Blob base structure
type Blob struct {
	ID       string `json:"id"`
	Size     int64  `json:"size"`
	RefCount int    `json:"refCount"`
	Data     []byte
}

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("blob.log", level)
}

func init() {
	InitLog(octlog.DebugLevel)
}

// DirPath to make blob path
func (b *Blob) DirPath() string {
	return utils.TrimDir(configuration.GetConfig().RootDirectory +
		manifest.BlobDir + "/" + b.ID[0:2] + "/" + b.ID[2:4])
}

// FilePath for blob
func (b *Blob) FilePath() string {
	return b.DirPath() + "/" + b.ID
}

// RefCountPath of blob
func (b *Blob) RefCountPath() string {
	return b.DirPath() + "/" + b.ID + ".refcount"
}

// GetRefCount of blob
func (b *Blob) GetRefCount() int {

	refcountFile := b.RefCountPath()
	if !utils.IsFileExist(refcountFile) {
		b.RefCount = 0
		return b.RefCount
	}

	data := utils.FileToString(refcountFile)
	if data == "" {
		utils.Remove(refcountFile)
		b.RefCount = 0
		return b.RefCount
	}

	b.RefCount = utils.StringToInt(data)
	return b.RefCount
}

// IsBlobExist for simple blob structure
func IsBlobExist(digest string) bool {

	b := Blob{
		ID: digest,
	}

	return utils.IsFileExist(b.FilePath())
}

// IsExist for simple blob structure
func (b *Blob) IsExist() bool {

	filepath := b.FilePath()
	if utils.IsFileExist(filepath) {
		return true
	}

	return false
}

// GetBlobPartial for partial blob fetching
func GetBlobPartial(name string, digest string) *Blob {
	b := Blob{
		ID: digest,
	}

	if !b.IsExist() {
		return nil
	}

	b.Size = utils.GetFileSize(b.FilePath())
	b.GetRefCount()

	return &b
}

// GetBlob to get blob from web api
func GetBlob(name string, digest string) *Blob {

	b := new(Blob)
	b.ID = digest

	if !b.IsExist() {
		octlog.Error("blob of %s not exist\n", digest)
		return nil
	}

	fd, err := os.Open(b.FilePath())
	if err != nil {
		octlog.Error("open file of %s error\n", b.FilePath())
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		octlog.Error("read file from %s error\n", b.FilePath())
		return nil
	}

	b.Data = data
	b.Size = utils.GetFileSize(b.FilePath())
	b.RefCount = b.GetRefCount()

	return b
}

// Delete to delete blob from api
func (b *Blob) Delete() error {
	if ref := b.DecRefCount(); ref == 0 {
		utils.Remove(b.FilePath())
		utils.Remove(b.RefCountPath())
	}
	return nil
}

// DecRefCount of blob
func (b *Blob) DecRefCount() int {
	if b.RefCount > 0 {
		b.RefCount = b.RefCount - 1
	}
	return 0
}

// IncRefCount of blob
func (b *Blob) IncRefCount() int {
	b.RefCount = b.RefCount + 1
	return b.RefCount
}

// WriteRefCount of blobs
func (b *Blob) WriteRefCount() {

	refcountFile := b.RefCountPath()
	if utils.IsFileExist(refcountFile) {
		utils.Remove(refcountFile)
	}

	fd, err := os.Create(refcountFile)
	if err != nil {
		octlog.Error("create file %s error %s\n", refcountFile, err)
	}

	defer fd.Close()

	_, err = fd.WriteString(utils.IntToString(b.RefCount))
	if err != nil {
		octlog.Warn("write refcount of %s error %s\n", b.ID, err)
	}
}

// WriteBlob To write blob from file
func (b *Blob) Write() error {

	utils.CreateDir(b.DirPath())

	// if blob exists, just increase refcount by 1
	if utils.IsFileExist(b.FilePath()) {
		b.IncRefCount()
		b.WriteRefCount()
		octlog.Warn("blob of %s already exist, just increase its refcount\n", b.ID)
		return nil
	}

	fd, err := os.Create(b.FilePath())
	if err != nil {
		octlog.Error("create blob of %s error\n", b.FilePath())
		return err
	}

	defer fd.Close()
	fd.Write(b.Data)

	b.IncRefCount()
	b.WriteRefCount()

	return nil
}

// ImportBlobs to write blobs from file and return its hash values
func ImportBlobs(filepath string) ([]string, int64, error) {

	f, err := os.Open(filepath)
	if err != nil {
		octlog.Error("file of %s not exist\n", filepath)
		return nil, 0, err
	}
	defer f.Close()

	var fileLength int64
	hashList := make([]string, 0)
	for {
		buffer := make([]byte, configuration.BlobSize)
		n, err := f.Read(buffer)
		if err == io.EOF {
			if n > 0 {
				dgst := utils.GetDigest(buffer[:n])
				octlog.Error("got size of %d,with hash:%s\n", n, dgst)
				b := &Blob{
					ID:   dgst,
					Data: buffer[:n],
				}
				b.GetRefCount()
				b.Write()

				//	WriteBlob(dgst, buffer[:n])
				fileLength += int64(n)
			}
			octlog.Warn("reached end of file[%d]\n", n)
			break
		}

		if err != nil {
			octlog.Error("read file error %s, %s bytes already read\n", err, fileLength)
			return nil, fileLength, err
		}

		fileLength += int64(n)
		dgst := utils.GetDigest(buffer[:n])
		b := &Blob{
			ID:   dgst,
			Data: buffer[:n],
		}
		b.GetRefCount()
		b.Write()

		//WriteBlob(dgst, buffer[:n])
		hashList = append(hashList, dgst)
	}

	return hashList, fileLength, nil
}

// HTTPGetBlob will get blob by name and digest
func HTTPGetBlob(url string) ([]byte, int, error) {

	resp, err := http.Get(url)
	if err != nil {
		octlog.Error("get url %s error\n", url)
		return nil, 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		octlog.Error("Read body from url %s error\n", url)
		return nil, 0, err
	}

	return body, int(resp.ContentLength), nil
}

// HTTPWriteBlob To write blob from file by HTTP
func HTTPWriteBlob(urlPattern string, dgst string, data []byte) error {

	url := urlPattern + "?digest=" + dgst
	reader := bytes.NewReader(data)
	reqeust, err := http.NewRequest("POST", url, reader)
	if err != nil {
		octlog.Error("New Http Request error on url %s\n", url)
		return err
	}

	resp, err := http.DefaultClient.Do(reqeust)
	if err != nil {
		octlog.Error("do http post error to url %s\n", url)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		octlog.Error("got bad status when post blob data %s\n", resp.Status)
		return errors.New("got bad status " + resp.Status)
	}

	octlog.Debug("HTTP upload blob %s to %s OK\n", dgst, url)

	return nil
}

// HTTPWriteBlobs to write blobs from file by HTTP
func HTTPWriteBlobs(filepath string, urlPattern string) ([]string, int64, error) {

	octlog.Debug("file %s, url %s\n", filepath, urlPattern)

	f, err := os.Open(filepath)
	if err != nil {
		octlog.Error("file of %s not exist\n", filepath)
		return nil, 0, err
	}
	defer f.Close()

	var fileLength int64
	hashList := make([]string, 0)
	for {
		buffer := make([]byte, configuration.BlobSize)
		n, err := f.Read(buffer)
		if err == io.EOF {
			octlog.Warn("reached end of file[%d]\n", n)
			break
		}
		fileLength += int64(n)

		if err != nil {
			octlog.Error("read file error %s", err)
			fmt.Printf("read file error %s\n", filepath)
			return hashList, fileLength, err
		}

		dgst := utils.GetDigest(buffer[:n])
		octlog.Debug("got size of %d,with hash:%s\n", n, dgst)

		err = HTTPWriteBlob(urlPattern, dgst, buffer[:n])
		if err != nil {
			octlog.Error("http post blob error url:%s,blob:%s\n", urlPattern, dgst)
			return hashList, fileLength, err
		}

		hashList = append(hashList, dgst)
	}

	octlog.Debug("file %s, url %s\n", filepath, urlPattern)

	return hashList, fileLength, nil
}
