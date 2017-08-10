package blobs

import (
	"fmt"
	"net/http"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func GetUploadStatus(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func StartBlobUpload(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func PatchBlobData(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func PutBlobUploadComplete(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func CancelBlobUpload(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}
