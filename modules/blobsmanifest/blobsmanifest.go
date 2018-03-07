package blobsmanifest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/config"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/octlog"
	"os"
)

var logger *octlog.LogConfig

// InitLog for blobs-manifest config
func InitLog(level int) {
	logger = octlog.InitLogConfig("blobs-manifest.log", level)
}

// BlobsManifest for base blobs and Manifest relationship
type BlobsManifest struct {
	Size    int64    `json:"size"`
	Chunks  []string `json:"chunks"`
	BlobSum string   `json:""`
}

// CalcBlobSum for blob chunks
func CalcBlobSum(chunks []string) string {
	var blobsum string
	for _, v := range chunks {
		blobsum += v
	}
	return utils.GetDigestStr(blobsum)
}

// GetBlobSum for sum hash value of blobs
func (bm *BlobsManifest) GetBlobSum() string {
	return CalcBlobSum(bm.Chunks)
}

func blobsManifestDirPath(dgst string) string {
	return utils.TrimDir(configuration.GetConfig().RootDirectory + manifest.BlobManifestDir + "/" + dgst[0:2])
}

func blobsManifestPath(dgst string) string {
	return blobsManifestDirPath(dgst) + "/" + dgst
}

// Delete blobs-manifest and blobs bellow it
func (bm *BlobsManifest) Delete() error {

	for _, chunk := range bm.Chunks {

		b := blobs.GetBlobPartial("", chunk)
		if b != nil {
			b.Delete()
		}
	}

	utils.Remove(blobsManifestPath(bm.BlobSum))

	return nil
}

// Write for manifest self delete
func (bm *BlobsManifest) Write() error {

	dirpath := blobsManifestDirPath(bm.BlobSum)
	utils.CreateDir(dirpath)

	filePath := dirpath + "/" + bm.BlobSum
	utils.Remove(filePath)

	fd, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("create file %s error\n", filePath)
		return err
	}

	defer fd.Close()

	data, _ := json.MarshalIndent(bm, "", "  ")
	fd.Write(data)

	return nil
}

// GetBlobsManifest to get blobs manifest config
func GetBlobsManifest(blobsum string) *BlobsManifest {

	bmPath := blobsManifestDirPath(blobsum) + "/" + blobsum
	if !utils.IsFileExist(bmPath) {
		logger.Errorf("file %s blobs-manifest not exist\n", blobsum)
		return nil
	}

	fd, err := os.Open(bmPath)
	if err != nil {
		logger.Errorf("open file of %s error %s\n", bmPath, err)
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		logger.Errorf("read data from %s error %s\n", bmPath, err)
	}

	bm := new(BlobsManifest)
	json.Unmarshal(data, bm)

	return bm
}

func readBlob(filepath string) []byte {

	fd, err := os.Open(filepath)
	if err != nil {
		logger.Errorf("open blob file %s error\n", filepath)
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		logger.Errorf("read data from blob file %s error\n", filepath)
		return nil
	}

	return data
}

// GetBlobHuge to get blob from huge file image.
func GetBlobHuge(blobSum string, digest string) *blobs.Blob {
	blobManifest := GetBlobsManifest(blobSum)
	if blobManifest == nil {
		logger.Errorf("blob-manifest %s not exist", blobSum)
		return nil
	}

	imageFilePath := configuration.RootDirectory() + "/" + manifest.ManifestDir + "/" + blobSum + "/" + "image"

	logger.Debugf("image file path of huge file %s", imageFilePath)

	b := new(blobs.Blob)
	b.ID = digest
	//b.Data = data
	//b.Size = utils.GetFileSize(b.FilePath())
	b.RefCount = 1

	return b	
}

// ExportHuge to export hugeblog file
func (bm *BlobsManifest)ExportHuge(outpath string, manifestDir string) error {
	if utils.IsFileExist(outpath) {
		os.Remove(outpath)
	}

	srcFilePath := manifestDir + "/" + bm.BlobSum + "/" + "image"
	_, err := utils.CopyFile(srcFilePath, outpath)

	return err
}

// Export file to outpath
func (bm *BlobsManifest)Export(outpath string) error {

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
		if hash == config.ZeroDataDigest8M {
			fd.Write(config.ZeroData8M)
			logger.Infof("no need to read zero data of %s\n", hash)
		} else {
			b := new(blobs.Blob)
			b.ID = hash
			data := readBlob(b.FilePath())
			fd.Write(data)
		}
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
