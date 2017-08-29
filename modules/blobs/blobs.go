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
	FilePath string `json:"filepath"`
	Size     int64  `json:"size"`
}

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("blob.log", level)
}

func init() {
	InitLog(octlog.DebugLevel)
}

// DirPath to make blob path
func DirPath(blobsum string) string {
	return utils.TrimDir(configuration.GetConfig().RootDirectory + manifest.BlobDir + "/" + blobsum[0:2] + "/" + blobsum[2:4])
}

// FilePath for blob
func FilePath(digest string) string {
	return DirPath(digest) + "/" + digest
}

// GetBlobSimple for simple blob structure
func GetBlobSimple(name string, digest string) *Blob {
	filepath := FilePath(digest)
	if utils.IsFileExist(filepath) {
		b := new(Blob)
		b.FilePath = filepath
		b.ID = digest
		return b
	}

	return nil
}

// GetBlob to get blob from web api
func GetBlob(name string, digest string) ([]byte, int, error) {

	blobpath := DirPath(digest) + "/" + digest
	if !utils.IsFileExist(blobpath) {
		octlog.Error("blob of %s not exist\n", blobpath)
		return nil, 0, errors.New("blob file of " + blobpath + " not exist")
	}

	fd, err := os.Open(blobpath)
	if err != nil {
		octlog.Error("open file of %s error\n", blobpath)
		return nil, 0, err
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		octlog.Error("read file from %s error\n", blobpath)
		return nil, 0, err
	}

	return data, int(utils.GetFileSize(blobpath)), nil
}

// DeleteBlob to delete blob from api
func DeleteBlob(name string, digest string) error {
	return nil
}

// GetUploadStatus to get upload status of blob
func GetUploadStatus(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// StartBlobUpload to start blob upload
func StartBlobUpload(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// PatchBlobData to patch blob data
func PatchBlobData(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// PutBlobUploadComplete to put blob upload complete
func PutBlobUploadComplete(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// CancelBlobUpload to cancel blob upload action
func CancelBlobUpload(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// WriteBlob To write blob from file
func WriteBlob(dgst string, data []byte) error {

	blobDir := DirPath(dgst)
	utils.CreateDir(blobDir)

	fd, err := os.Create(blobDir + "/" + dgst)
	if err != nil {
		octlog.Error("create blob of %s error\n", blobDir+"/"+dgst)
		return err
	}

	defer fd.Close()
	fd.Write(data)

	return nil
}

// WriteBlobs to write blobs from file and return its hash values
func WriteBlobs(filepath string) ([]string, int64, error) {

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
			octlog.Error("reached end of file[%d]\n", n)
			break
		}
		fileLength += int64(n)

		if err != nil {
			octlog.Error("read file error %s", err)
		}

		dgst := utils.GetDigest(buffer[:n])
		octlog.Error("got size of %d,with hash:%s\n", n, dgst)
		WriteBlob(dgst, buffer[:n])

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

	octlog.Error("HTTP upload blob %s to %s OK\n", dgst, url)

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
			octlog.Error("reached end of file[%d]\n", n)
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

	fmt.Printf("file %s, url %s\n", filepath, urlPattern)

	return hashList, fileLength, nil
}
