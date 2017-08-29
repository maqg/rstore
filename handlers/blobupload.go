package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobupload"
	"octlink/rstore/utils"
	"octlink/rstore/utils/serviceresp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func blobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	resp, err := blobupload.UploadBlob(name, digest)
	if err != nil {
		serviceresp.NotFoundResp(w)
		return
	}

	data := utils.JSON2String(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

// if digest specified,
func startBlobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	resp, err := blobupload.UploadBlob(name, digest)
	if err != nil {
		serviceresp.NotFoundResp(w)
		return
	}

	data := utils.JSON2String(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

func deleteBlobUpload(w http.ResponseWriter, r *http.Request) {

}

func blobUploadManager(r *http.Request) http.Handler {
	return handlers.MethodHandler{
		"POST":   http.HandlerFunc(startBlobUpload),
		"PATCH":  http.HandlerFunc(blobUpload),
		"DELETE": http.HandlerFunc(deleteBlobUpload),
	}
}
