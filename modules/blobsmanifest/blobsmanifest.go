package blobsmanifest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
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
		octlog.Error("create file %s error\n", filePath)
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

	if !utils.IsFileExist(bmPath) {
		octlog.Error("file %s blobs-manifest not exist\n", blobsum)
		return nil
	}

	fd, err := os.Open(bmPath)
	if err != nil {
		octlog.Error("open file of %s error %s\n", bmPath, err)
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		octlog.Error("read data from %s error %s\n", bmPath, err)
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

// HTTPGetBlobsManifest will get blobs manifest by name and digest
func HTTPGetBlobsManifest(url string) (*BlobsManifest, error) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("get url %s error\n", url)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body from url %s error\n", url)
		return nil, err
	}

	bm := new(BlobsManifest)
	err = json.Unmarshal(body, bm)
	if err != nil {
		fmt.Printf("parse body to blobsmanifest error[%s]\n", string(body))
		return nil, err
	}

	return bm, nil
}

// HTTPWrite for manifest to write by http
func (bm *BlobsManifest) HTTPWrite(url string) error {

	data, err := json.Marshal(bm.Chunks)
	if err != nil {
		octlog.Error("convert chunks to json bytes error\n")
		return err
	}

	url += fmt.Sprintf("?size=%d", bm.Size)
	reqeust, err := http.NewRequest("POST", url, bytes.NewReader(data))
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
		octlog.Error("got bad status when post blob data %s,url:%s\n", resp.Status, url)
		return errors.New("got bad status " + resp.Status)
	}

	octlog.Debug("HTTP upload blobsmanifest %s to %s OK\n", bm.BlobSum, url)

	return nil
}
