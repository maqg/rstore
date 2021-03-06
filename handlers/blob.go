package handlers

import (
	"octlink/rstore/modules/blobsmanifest"
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobs"
	"octlink/rstore/utils"
	"octlink/rstore/utils/serviceresp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// RenderErrorMsg to render error msg
func RenderErrorMsg(w http.ResponseWriter, r *http.Request, errMsg string) {
	msg := utils.JSON2String(serviceresp.SuccessResp(errMsg))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(msg)))
	fmt.Fprint(w, msg)
}

// RenderMsg for render msg for success
func RenderMsg(w http.ResponseWriter, r *http.Request, data interface{}) {
	msg := utils.JSON2String(serviceresp.SuccessResp(data))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(msg)))
	fmt.Fprint(w, msg)
}

// GetBlob to get blob from web api
func getBlob(w http.ResponseWriter, r *http.Request) {

	var b *blobs.Blob

	name := mux.Vars(r)["name"] // now means blobsum
	digest := mux.Vars(r)["digest"]

	dgst, index, length := utils.ParseBlobDigest(digest)
	if length != 0 { // get blob digest like
		b = blobsmanifest.GetBlobHuge(name, dgst, index, length)
	} else {
		b = blobs.GetBlob(name, digest)
	}
	if b == nil {
		logger.Errorf("get blob by %s:%s error\n", name, digest)
		RenderErrorMsg(w, r, "blob of "+digest+" not exist")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", b.Size))
	n, err := w.Write(b.Data)
	if err != nil {
		logger.Errorf("Write blob %s to client error\n", digest)
		return
	}

	logger.Debugf("Write to Client %d bytes blob %s OK\n", n, digest)
}

// DeleteBlob to delete blob from api
func deleteBlob(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	b := blobs.GetBlobPartial(name, digest)
	if b == nil {
		w.WriteHeader(http.StatusNotFound)
		logger.Errorf("blob digest of %s not exist to delete\n", digest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func blobManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(getBlob),
		"HEAD": http.HandlerFunc(getBlob),
	}

	mhandler["DELETE"] = http.HandlerFunc(deleteBlob)

	return mhandler
}
