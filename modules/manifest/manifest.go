package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/rstore/configuration"
	"octlink/rstore/utils"
	"octlink/rstore/utils/octlog"
	"os"
)

// Manifest base Manifest structure
type Manifest struct {
	ID          string `json:"uuid"`
	Name        string `json:"name"`    // Image UUID
	BlobSum     string `json:"blobsum"` // Sum of Blob Digests
	DiskSize    int64  `json:"diskSize"`
	VirtualSize int64  `json:"virtualSize"`
	CreateTime  string `json:"createTime"`
	Path        string
}

const (
	// ReposDir Base Repos Directory
	ReposDir = "/registry/repos"

	// BlobManifestDir for manifest-blobs relationship
	BlobManifestDir = "/registry/blob-manifests"

	// BlobDir for blobs tree directory
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
func GetManifest(name string, dgst string) *Manifest {

	octlog.Debug("got name[%s],digest[%s]\n", name, dgst)

	conf := configuration.GetConfig()
	maniPath := conf.RootDirectory + fmt.Sprintf(ManifestFileProto, name, dgst)
	if !utils.IsFileExist(maniPath) {
		octlog.Error("manifest not exist %s\n", maniPath)
		logger.Errorf("manifest not exist %s\n", maniPath)
		return nil
	}

	manifest, err := FileToManifest(maniPath)
	if err != nil {
		octlog.Error("manifest parse error %s\n", maniPath)
		logger.Errorf("manifest parse error %s\n", maniPath)
		return nil
	}

	return manifest
}

func init() {
	InitLog(octlog.DEBUG_LEVEL)
}

// Delete for manifest self delete
func (manifest *Manifest) Delete() error {
	return nil
}

// Write for manifest self delete
func (manifest *Manifest) Write() error {

	// create manifest base diretory
	utils.CreateDir(manifest.Path)

	filePath := manifest.Path + fmt.Sprintf("/%s.json", manifest.ID)
	utils.Remove(filePath)

	fd, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create file %s error\n", filePath)
		return err
	}

	defer fd.Close()

	data, _ := json.MarshalIndent(manifest, "", "  ")
	fd.Write(data)

	return nil
}
