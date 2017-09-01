package api

import (
	"octlink/rstore/modules/task"
)

// ShowAllTasks by condition
func ShowAllTasks(paras *Paras) *Response {

	resp := new(Response)
	resp.Data = task.GetTasks()

	return resp
}

// ShowTask get one task by task id
func ShowTask(paras *Paras) *Response {
	resp := new(Response)
	return resp
}

// DeleteTask task by id
func DeleteTask(paras *Paras) *Response {
	resp := new(Response)
	return resp
}

// StopTask by id
func StopTask(paras *Paras) *Response {
	resp := new(Response)
	return resp
}
