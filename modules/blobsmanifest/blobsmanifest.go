package blobsmanifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"octlink/rstore/configuration"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"os"
)

// BlobsManifest for base blobs and Manifest relationship
type BlobsManifest struct {
	Size    int64    `json:"size"`
	Chunks  []string `json:"chunks"`
	BlobSum string   `json:""`
}

// GetBlobSum for sum hash value of blobs
func (bm *BlobsManifest) GetBlobSum() string {
	var blobsum string
	for _, v := range bm.Chunks {
		blobsum += v
	}
	return utils.GetDigestStr(blobsum)
}

func dirpath(dgst string) string {
	return utils.TrimDir(configuration.GetConfig().RootDirectory + manifest.BlobManifestDir + "/" + dgst[0:2])
}

// Write for manifest self delete
func (bm *BlobsManifest) Write() error {

	dirpath := dirpath(bm.BlobSum)
	utils.CreateDir(dirpath)

	filePath := dirpath + "/" + bm.BlobSum
	utils.Remove(filePath)

	fd, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("create file %s error\n", filePath)
		return err
	}

	defer fd.Close()

	data, _ := json.MarshalIndent(bm, "", "  ")
	fd.Write(data)

	return nil
}

// GetBlobsManifest to get blobs manifest config
func GetBlobsManifest(blobsum string) *BlobsManifest {
	bmPath := dirpath(blobsum) + "/" + blobsum
	fd, err := os.Open(bmPath)
	if err != nil {
		fmt.Printf("open file of %s error %s\n", bmPath, err)
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Printf("read data from %s error %s\n", bmPath, err)
	}

	bm := new(BlobsManifest)
	json.Unmarshal(data, bm)

	return bm
}

func readBlob(filepath string) []byte {
	fd, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open blob file %s error\n", filepath)
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Printf("read data from blob file %s error\n", filepath)
		return nil
	}

	return data
}

// Export file to outpath
func (bm *BlobsManifest) Export(outpath string) error {

	if utils.IsFileExist(outpath) {
		os.Remove(outpath)
	}

	fd, err := os.Create(outpath)
	if err != nil {
		fmt.Printf("create file of %s error\n", outpath)
		return err
	}

	defer fd.Close()

	for _, hash := range bm.Chunks {
		blobPath := blobs.DirPath(hash) + "/" + hash
		data := readBlob(blobPath)
		fd.Write(data)
	}

	return nil
}
