package serviceresp

import (
	"net/http"
)

// ServiceResp for Service Responce base structure
type ServiceResp struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	ErrMsg string      `json:"errMsg"`
}

// SuccessResp build and returen successfull response
func SuccessResp(data interface{}) *ServiceResp {
	sr := new(ServiceResp)
	sr.Data = data
	sr.Status = "ok"
	return sr
}

// FailedResp build and return failed response
func FailedResp(errMsg string) *ServiceResp {
	sr := new(ServiceResp)
	sr.Status = "error"
	sr.ErrMsg = errMsg
	return sr
}

// NotFoundResp for not found http response
func NotFoundResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

// StatusOKResp for status ok response
func StatusOKResp(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
