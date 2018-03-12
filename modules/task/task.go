package task

import (
	"fmt"
	"io"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
	"octlink/rstore/utils/merrors"
	"octlink/rstore/utils/octlog"
	"octlink/rstore/utils/uuid"
	"os"
	"strings"
)

const (
	// TaskStatusNew to new status task
	TaskStatusNew = "new"

	// TaskStatusRunning for running status task
	TaskStatusRunning = "running"

	// TaskStatusFinished for finished status task
	TaskStatusFinished = "finished"

	// TaskStatusError for error status task
	TaskStatusError = "error"

	// TaskStatusStopped for stopped status task
	TaskStatusStopped = "stopped"
)

// Task for downloading and importing task management
type Task struct {
	ID         string `json:"id"`
	URL        string `josn:"url"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
	FinishTime string `json:"finishTime"`
	FilePath   string `json:"filePath"`
	FileName   string `json:"fileName"`
	FileLength int64  `json:"fileLength"`
	ImageName  string `json:"imageName"`
	callback   ImageCallBack
}

var logger *octlog.LogConfig

// InitLog for tasks
func InitLog(level int) {
	logger = octlog.InitLogConfig("task.log", level)
}

// GTasks all tasks map management
var GTasks = make(map[string]*Task, 0)

// GetTasks from web api
func GetTasks() map[string]*Task {
	return GTasks
}

// GetTaskByImage will get task by image id
func GetTaskByImage(imageID string) *Task {
	return nil
}

// ImageCallBack for image callback
type ImageCallBack func(string, int64, int64, string, string) error

// AddAndRun will add a new task to GTasks and run it
func (t *Task) AddAndRun(callback ImageCallBack) {
	t.callback = callback
	t.Status = TaskStatusRunning

	GTasks[t.ID] = t
	go t.Download()

	logger.Warnf("task of %s start to run, %s\n", t.ID, t.URL)
}

// GetTask by taskid
func GetTask(id string) *Task {
	return GTasks[id]
}

// UpdateFilePath update and write file path
func (t *Task) UpdateFilePath() {

	if t.URL[:4] == "http" || t.URL[:4] == "ftp" {
		segs := strings.Split(t.URL, "/")
		t.FileName = segs[len(segs)-1]
		t.FilePath = configuration.GetConfig().RootDirectory + "/registry/temp/" +
				uuid.Generate().Simple() + "EEEEE" + t.FileName		
	} else { // for local file import
		segs := strings.Split(t.URL, "/")
		t.FileName = segs[len(segs)-1]
		t.FilePath = t.URL
	}

	logger.Debugf("update file dst path %s\n", t.FilePath)
}

// ImportBlobs and then write blobs-manifest config
func (t *Task) ImportBlobs() (*blobsmanifest.BlobsManifest, error) {

	var hashes []string
	var len int64
	var err error

	hugeBlob := configuration.GetConfig().HugeBlob
	if hugeBlob {
		hashes, len, err = blobs.ImportHugeBlob(t.FilePath)
	} else {
		hashes, len, err = blobs.ImportBlobs(t.FilePath)
	}
	if err != nil {
		logger.Errorf("got file hashlist error\n")
		return nil, err
	}

	if len != t.FileLength {
		logger.Errorf("filelen of blobs and http contentlen not match %d:%d\n", len, t.FileLength)
		return nil, fmt.Errorf("filelen %d not match imported len of %d", t.FileLength, len)
	}

	blobsum := blobsmanifest.CalcBlobSum(hashes)
	bm := blobsmanifest.GetBlobsManifest(blobsum)
	if bm != nil {
		logger.Warnf("blobs-manifest %s already exist, just return it\n", bm.BlobSum)
		return bm, nil
	}

	// write blobs-manifest config
	bm = new(blobsmanifest.BlobsManifest)
	bm.Size = t.FileLength
	bm.Chunks = hashes
	bm.BlobSum = blobsum
	err = bm.Write()
	if err != nil {
		logger.Errorf("write blobs-manifest %s error\n", bm.BlobSum)
		return nil, err
	}

	return bm, nil
}

// WriteManifest to write manifest file
func (t *Task) WriteManifest(blobsum string) (*manifest.Manifest, error) {

	m := manifest.GetManifest(t.ImageName, blobsum)
	if m != nil {
		logger.Warnf("manifest of %s:%s already exist\n", t.ImageName, blobsum)
		return m, nil
	}

	// write manifest config
	m = new(manifest.Manifest)
	m.Name = t.ImageName
	m.DiskSize = t.FileLength
	m.VirtualSize = utils.GetVirtualSize(t.FilePath)
	m.CreateTime = utils.CurrentTimeStr()
	m.BlobSum = blobsum

	err := m.Write()
	if err != nil {
		logger.Errorf("Create manifest error[%s]\n", err)
		return nil, err
	}

	logger.Debugf("Write manifest of %s OK\n", blobsum)

	return m, nil
}


func importImage(t *Task) {
	bm, err := t.ImportBlobs()
	if err != nil {
		logger.Errorf("import blobs error %s for %s\n", err, t.FileName)
		t.Error()

		// Remove temp file when failed to add image
		if t.URL[:4] == "http" || t.URL[:4] == "ftp" {
			utils.Remove(t.FilePath)
		}
		return
	}

	m, err := t.WriteManifest(bm.BlobSum)
	if err != nil {
		logger.Errorf("write manifest error, imageid:%s, blobsum:%s\n", t.ImageName, bm.BlobSum)
		t.Error()
		// Remove temp file when failed to add image
		if t.URL[:4] == "http" || t.URL[:4] == "ftp" {
			utils.Remove(t.FilePath)
		}
		return
	}

	hugeBlob := configuration.GetConfig().HugeBlob
	if hugeBlob {
		if utils.IsFileExist(t.FilePath) {
			//
		}
	}

	// update task status and image info
	t.Finish(m.DiskSize, m.VirtualSize, m.BlobSum)

	logger.Debugf("got file length of %d, and wroted %s\n", t.FileLength, t.FilePath)
}

// Download Image from URL
func (t *Task) Download() {

	if t.URL[:4] != "http" && t.URL[:4] != "ftp" {

		t.UpdateFilePath()

		// import image from local
		if !utils.IsFileExist(t.URL) {
			logger.Errorf("local url to import of %s not exist", t.URL)
			t.Error()
			return
		}
		t.FileLength = utils.GetFileSize(t.URL)

		importImage(t)
		return
	}

	r, err := http.Get(t.URL)
	if err != nil {
		logger.Errorf("create URL response error on url %s\n", t.URL)
		t.Error()
		return
	}
	defer r.Body.Close()

	t.UpdateFilePath()
	fd, err := os.Create(t.FilePath)
	if err != nil {
		logger.Errorf("create temp file %s error\n", t.FilePath)
		t.Error()
		return
	}

	defer fd.Close()

	buf := make([]byte, configuration.BlobSize)
	for {

		n, err := r.Body.Read(buf)
		if err == io.EOF {
			if n > 0 {
				t.FileLength += int64(n)
				wlen, err := fd.Write(buf[:n])
				if err != nil {
					logger.Warnf("got len %d and wroted %d, total %d\n", n, wlen, t.FileLength)
					t.Error()
					return
				}
			}
			break
		}

		if err != nil {
			logger.Errorf("read file error %s", err)
			// TBD clear
			t.Error()
			return
		}

		t.FileLength += int64(n)
		wlen, err := fd.Write(buf[:n])
		if err != nil {
			logger.Errorf("got len %d and wroted %d, total %d\n", n, wlen, t.FileLength)
			t.Error()
			return
		}
	}

	importImage(t)

	return
}

// Stop this task
func (t *Task) Stop() int {
	t.callback(t.ImageName, 0, 0, "", TaskStatusError)
	return merrors.ErrSuccess
}

// Delete this task
func (t *Task) Delete() int {
	return merrors.ErrSuccess
}

// Finish task here
func (t *Task) Finish(diskSize int64, virtualSize int64, blobsum string) {
	t.Status = TaskStatusFinished
	t.FinishTime = utils.CurrentTimeStr()
	t.callback(t.ImageName, diskSize, virtualSize, blobsum, TaskStatusFinished)
}

func (t *Task) Error() {
	t.Status = TaskStatusError
	t.FinishTime = utils.CurrentTimeStr()
	t.callback(t.ImageName, 0, 0, "", TaskStatusError)
}
