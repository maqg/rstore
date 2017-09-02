package task

import (
	"io"
	"net/http"
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
type ImageCallBack func(string, int64)

// AddAndRun will add a new task to GTasks and run it
func (t *Task) AddAndRun(callback ImageCallBack) {
	GTasks[t.ID] = t
	t.Run()
	callback(t.ImageName, 9000)
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

	// update task status and image info
	t.Finish()
	octlog.Warn("got file length of %d, and wroted %s\n", t.FileLength, t.FilePath)

	return
}

// Run this task
func (t *Task) Run() {
	t.Status = TaskStatusRunning
	go t.Download()
	octlog.Warn("task of %s start to run, %s\n", t.ID, t.URL)
}

// Stop this task
func (t *Task) Stop() error {
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
}

func (t *Task) Error() {
	t.Status = TaskStatusError
	t.FinishTime = utils.CurrentTimeStr()
}
