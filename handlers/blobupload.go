package handlers

import (
	"octlink/rstore/utils/configuration"
	"octlink/rstore/modules/manifest"
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/modules/blobupload"
	"octlink/rstore/utils/serviceresp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func blobHugeUpload(w http.ResponseWriter, r *http.Request, name string, digest string, tempName string) {
	tempFile := configuration.RootDirectory() + manifest.TempDir + "/" + tempName
	err := blobupload.CopyFullPayload(w, r, tempFile)
	if err != nil {
		logger.Errorf("copy full data for blob %s error\n", digest)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	serviceresp.StatusOKResp(w)
}

func blobUpload(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")
	tempName := r.FormValue("tempName")

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

	if tempName != "" {
		fmt.Printf("running in blob huge uploading %s,digest %s\n", tempName, digest)
		blobHugeUpload(w, r, name, digest, tempName)
		return
	}
	fmt.Printf("running in blob uploading %s,digest %s\n", tempName, digest)

	b := blobs.Blob{
		ID: digest,
	}

	if b.IsExist() {
		b.IncRefCount()
		b.WriteRefCount()
		serviceresp.StatusOKResp(w)
		logger.Warnf("blob of %s already exist", digest)
		return
	}

	logger.Debugf("start to upload blob %s\n", digest)

	bu := blobupload.BlobUpload{
		ID:         digest,
		FilePath:   b.FilePath(),
		RespWriter: w,
		Request:    r,
	}
	err := bu.Upload()
	if err != nil {
		logger.Errorf("upload blob of %s failed\n", bu.FilePath)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	logger.Debugf("Upload blob %s OK\n", digest)

	serviceresp.StatusOKResp(w)
}

func blobUploadManager(r *http.Request) http.Handler {
	return handlers.MethodHandler{
		"POST":  http.HandlerFunc(blobUpload),
		"PATCH": http.HandlerFunc(blobUpload),
	}
}
