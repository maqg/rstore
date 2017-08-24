package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/configuration"
	"octlink/rstore/utils"
	"octlink/rstore/utils/octlog"
	"os"

	"github.com/gorilla/mux"
)

// Manifest base Manifest structure
type Manifest struct {
	ID          string `json:"uuid"`
	Name        string `json:"name"`    // Image UUID
	BlobSum     string `json:"blobsum"` // Sum of Blob Digests
	DiskSize    int64  `json:"diskSize"`
	VirtualSize int64  `json:"virtualSize"`
	CreateTime  string `json:"createTime"`
}

const (
	// ReposDir Base Repos Directory
	ReposDir = "/registry/repos"

	// BlobManifestDir for manifest-blobs relationship
	BlobManifestDir = "/registry/blob-manifests"

	// BlobsDir for blobs tree directory
	BlobDir = "/registry/blobs"

	// ImageDirProto Image Dir Proto Type
	ImageDirProto = "/registry/repos/%s"

	// ManifestDirProto manifest directory proto type
	ManifestDirProto = "/registry/repos/%s/manifests"

	// ManifestFileProto manifest file of json proto type
	ManifestFileProto = "/registry/repos/%s/manifests/%s.json"
)

var logger *octlog.LogConfig

// InitLog for manifest module
func InitLog(level int) {
	logger = octlog.InitLogConfig("manifest.log", level)
}

// FileToManifest load file to Manifest struct from json
func FileToManifest(filePath string) (*Manifest, error) {

	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	in, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	manifest := new(Manifest)
	if err := json.Unmarshal(in, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

// GetManifest for api call
func GetManifest(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]

	octlog.Debug("got name[%s],digest[%s]\n", name, digest)

	conf := configuration.GetConfig()
	maniPath := conf.RootDirectory + fmt.Sprintf(ManifestFileProto, name, digest)
	if !utils.IsFileExist(maniPath) {
		w.WriteHeader(http.StatusNotFound)
		octlog.Error("manifest not exist %s\n", maniPath)
		logger.Errorf("manifest not exist %s\n", maniPath)
		return
	}

	manifest, err := FileToManifest(maniPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		octlog.Error("manifest parse error %s\n", maniPath)
		logger.Errorf("manifest parse error %s\n", maniPath)
		return
	}

	data, _ := json.Marshal(manifest)
	dataStr := utils.BytesToString(data)

	octlog.Debug("Got manifest %s\n", dataStr)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(dataStr)))
	fmt.Fprint(w, dataStr)
}

// DeleteManifest for api call
func DeleteManifest(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// PutManifest for api call
func PutManifest(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func init() {
	InitLog(octlog.DEBUG_LEVEL)
}
