package task

import (
	"io"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/configuration"
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
	Callback   ImageCallBack
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
	GTasks[t.ID] = t
	t.Run(callback)
}

// GetTask by taskid
func GetTask(id string) *Task {
	return nil
}

// UpdateFilePath update and write file path
func (t *Task) UpdateFilePath() {
	segs := strings.Split(t.URL, "/")
	t.FileName = segs[len(segs)-1]
	t.FilePath = configuration.GetConfig().RootDirectory + "/registry/temp/" +
		uuid.Generate().Simple() + "EEEEE" + t.FileName
}

// Download Image from URL
func (t *Task) Download() {

	r, err := http.Get(t.URL)
	if err != nil {
		octlog.Error("get url %s error\n", t.URL)
		t.Error()
		return
	}

	defer r.Body.Close()

	t.UpdateFilePath()
	fd, err := os.Create(t.FilePath)
	if err != nil {
		octlog.Error("create temp file %s error\n", t.FilePath)
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
					octlog.Warn("got len %d and wroted %d, total %d\n", n, wlen, t.FileLength)
					t.Error()
					return
				}
			}
			break
		}

		if err != nil {
			octlog.Error("read file error %s", err)
			// TBD clear
			t.Error()
			return
		}

		t.FileLength += int64(n)
		wlen, err := fd.Write(buf[:n])
		if err != nil {
			octlog.Warn("got len %d and wroted %d, total %d\n", n, wlen, t.FileLength)
			t.Error()
			return
		}
	}

	hashes, len, err := blobs.ImportBlobs(t.FilePath)
	if err != nil {
		octlog.Error("got file hashlist error\n")
		return
	}

	if len != t.FileLength {
		octlog.Error("filelen of blobs and http contentlen not match %d:%d\n", len, t.FileLength)
		return
	}

	// write blobs-manifest config
	bm := new(blobsmanifest.BlobsManifest)
	bm.Size = t.FileLength
	bm.Chunks = hashes
	bm.BlobSum = bm.GetBlobSum()
	err = bm.Write()
	if err != nil {
		octlog.Error("write blobs-manifest error\n")
		return
	}

	// write manifest config
	mid := utils.GetDigestStr(t.ImageName)
	manifest := new(manifest.Manifest)
	manifest.Name = t.ImageName
	manifest.ID = mid
	manifest.DiskSize = t.FileLength
	manifest.VirtualSize = utils.GetVirtualSize(t.FilePath)
	manifest.CreateTime = utils.CurrentTimeStr()
	manifest.BlobSum = bm.BlobSum

	err = manifest.Write()
	if err != nil {
		octlog.Error("Create manifest error[%s]\n", err)
		// TDB,rollback
		return
	}

	// update task status and image info
	t.Finish()
	octlog.Warn("got file length of %d, and wroted %s\n", t.FileLength, t.FilePath)

	return
}

// Run this task
func (t *Task) Run(callback ImageCallBack) {
	t.Status = TaskStatusRunning
	go t.Download()
	octlog.Warn("task of %s start to run, %s\n", t.ID, t.URL)
}

// Stop this task
func (t *Task) Stop(callback ImageCallBack) error {
	callback(t.ImageName, 0, 0, "", TaskStatusError)
	return nil
}

// Delete this task
func (t *Task) Delete() error {
	return nil
}

// Finish task here
func (t *Task) Finish() {
	t.Status = TaskStatusFinished
	t.FinishTime = utils.CurrentTimeStr()
	t.Callback(t.ImageName, 0, 0, "", TaskStatusError)
}

func (t *Task) Error() {
	t.Status = TaskStatusError
	t.FinishTime = utils.CurrentTimeStr()
	t.Callback(t.ImageName, 0, 0, "", TaskStatusError)
}
