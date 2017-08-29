package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobupload"
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
		serviceresp.StatusOKResp(w)
		fmt.Printf("blob of %s already exist", digest)
		return
	}

	fmt.Printf("start to upload blob %s\n", digest)

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
	w.WriteHeader(http.StatusOK)
}

func deleteBlobUpload(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func blobUploadManager(r *http.Request) http.Handler {
	return handlers.MethodHandler{
		"POST":   http.HandlerFunc(blobUpload),
		"PATCH":  http.HandlerFunc(blobUpload),
		"DELETE": http.HandlerFunc(deleteBlobUpload),
	}
}
