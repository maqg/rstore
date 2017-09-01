package task

import (
	"fmt"
	"io"
	"net/http"
	"octlink/rstore/utils/octlog"
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
	ImageName  string `json:"imageName"`
}

// GTasks all tasks map management
var GTasks = make(map[string]*Task, 0)

// GetTaskByImage will get task by image id
func GetTaskByImage(imageID string) *Task {
	return nil
}

// Add will add a new task to GTasks and run it
func (t *Task) Add() error {

	GTasks[t.ID] = t
	t.Status = TaskStatusRunning
	t.Run()

	return nil
}

// GetTask by taskid
func GetTask(id string) *Task {
	return nil
}

// GetTasks will return all running tasks
func GetTasks() []*Task {
	return nil
}

// Download Image from URL
func (t *Task) Download() error {

	r, err := http.Get(t.URL)
	if err != nil {
		fmt.Printf("get url %s error\n", t.URL)
		return nil
	}

	defer r.Body.Close()

	var fileLength int64

	buf := make([]byte, 1024*4)
	for {
		n, err := r.Body.Read(buf)
		if err == io.EOF {
			octlog.Warn("got len %d\n", n)
			fileLength += int64(n)
			break
		}

		octlog.Warn("got new len %d\n", n)
		fileLength += int64(n)

		if err != nil {
			octlog.Error("read file error %s", err)
			// TBD clear
			return nil
		}

		//buffer[:n]
	}

	// update task status and image info
	octlog.Warn("got file length of %d\n", fileLength)

	return nil
}

// Run this task
func (t *Task) Run() error {
	go t.Download()
	return nil
}

// Stop this task
func (t *Task) Stop() error {
	return nil
}

// Delete this task
func (t *Task) Delete() error {
	return nil
}
