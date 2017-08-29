package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobupload"
	"octlink/rstore/utils"
	"octlink/rstore/utils/serviceresp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func blobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	if name == "" || digest == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "" && ct != "application/octet-stream" {
		fmt.Printf("Bad Content-Type\n")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if b := blobs.GetBlobSimple(name, digest); b != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("blob of %s already exist", digest)
		return
	}

	bu := blobupload.BlobUpload{
		ID:         digest,
		FilePath:   blobs.FilePath(digest),
		RespWriter: w,
		Request:    r,
	}
	err := bu.Upload()
	if err != nil {
		fmt.Printf("upload blob of %s failed\n", bu.FilePath)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	serviceresp.StatusOKResp(w)
}

// if digest specified,
func startBlobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	resp, err := blobupload.StartBlobUpload(name, digest)
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
