package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobupload"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startBlobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	resp, err := blobupload.StartBlobUpload(name, digest)
	if err != nil {
		writeBackError(w, r)
		return
	}

	data := utils.JSON2String(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

func patchBlobData(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func putBlobUploadComplete(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func cancelBlobUpload(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// GetUploadStatus to get upload status of blob
func getUploadStatus(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func blobUploadManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(getUploadStatus),
		"HEAD": http.HandlerFunc(getUploadStatus),
	}

	mhandler["POST"] = http.HandlerFunc(startBlobUpload)
	mhandler["PATCH"] = http.HandlerFunc(patchBlobData)
	mhandler["PUT"] = http.HandlerFunc(putBlobUploadComplete)
	mhandler["DELETE"] = http.HandlerFunc(cancelBlobUpload)

	return mhandler
}
