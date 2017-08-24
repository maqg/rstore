package blobsmanifest

import (
	"encoding/json"
	"fmt"
	"octlink/rstore/utils"
	"os"
)

// BlobsManifest for base blobs and Manifest relationship
type BlobsManifest struct {
	Size    int64    `json:"size"`
	Chunks  []string `json:"chunks"`
	BlobSum string   `json:""`
	Path    string
}

// GetBlobSum for sum hash value of blobs
func (bm *BlobsManifest) GetBlobSum() string {
	var blobsum string
	for _, v := range bm.Chunks {
		blobsum += v
	}
	return utils.GetDigestStr(blobsum)
}

func (bm *BlobsManifest) dirpath() string {
	return bm.Path + "/" + bm.BlobSum[0:2]
}

// Write for manifest self delete
func (bm *BlobsManifest) Write() error {

	dirpath := bm.dirpath()
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
