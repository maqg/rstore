package api

import (
	"octlink/rstore/modules/task"
)

// ShowAllTasks by condition
func ShowAllTasks(paras *Paras) *Response {
	return &Response{
		Data: task.GetTasks(),
	}
}

// ShowTask get one task by task id
func ShowTask(paras *Paras) *Response {
	return &Response{
		Data: task.GetTask(paras.Get("id")),
	}
}

// DeleteTask task by id
func DeleteTask(paras *Paras) *Response {
	t := task.GetTask(paras.Get("id"))
	return &Response{
		Error: t.Delete(),
	}
}

// StopTask by id
func StopTask(paras *Paras) *Response {
	t := task.GetTask(paras.Get("id"))
	return &Response{
		Error: t.Stop(),
	}
}
