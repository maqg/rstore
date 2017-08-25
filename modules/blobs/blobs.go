package blobs

import (
	"fmt"
	"io"
	"net/http"
	"octlink/rstore/configuration"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/octlog"
	"os"

	"github.com/gorilla/mux"
)

var logger *octlog.LogConfig

// InitLog to init log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("blob.log", level)
}

// GetBlob to get blob from web api
func GetBlob(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	emptyJSON := fmt.Sprintf("{\"msg\":\"this is blob message,name:%s,digest:%s\"}", name, digest)

	logger.Debugf("got name:%s,digest:%s\n", name, digest)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// DeleteBlob to delete blob from api
func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
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

// GetBlobDir to get blob dir from root and digest
func GetBlobDir(rootdirectory string, digest string) string {
	return rootdirectory + manifest.BlobDir + "/" + digest[0:2] + "/" + digest[2:4]
}

// WriteBlob To write blob from file
func WriteBlob(rootdirectory string, dgst string, data []byte) error {

	blobDir := GetBlobDir(rootdirectory, dgst)
	utils.CreateDir(blobDir)

	fd, err := os.Create(blobDir + "/" + dgst)
	if err != nil {
		fmt.Printf("create blob of %s error\n", blobDir+"/"+dgst)
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
		fmt.Printf("file of %s not exist\n", filepath)
		return nil, 0, err
	}
	defer f.Close()

	var fileLength int64
	hashList := make([]string, 0)
	for {
		buffer := make([]byte, configuration.BlobSize)
		n, err := f.Read(buffer)
		if err == io.EOF {
			fmt.Printf("reached end of file[%d]\n", n)
			break
		}
		fileLength += int64(n)

		if err != nil {
			fmt.Printf("read file error %s", err)
		}

		dgst := utils.GetDigest(buffer[:n])
		fmt.Printf("got size of %d,with hash:%s\n", n, dgst)
		WriteBlob(configuration.GetConfig().RootDirectory, dgst, buffer[:n])

		hashList = append(hashList, dgst)
	}

	return hashList, fileLength, nil
}

// DirPath to make blob path
func DirPath(blobsum string) string {
	return utils.TrimDir(configuration.GetConfig().RootDirectory + manifest.BlobDir + "/" + blobsum[0:2] + "/" + blobsum[2:4])
}
