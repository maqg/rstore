package handlers

import (
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

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]
	// digest := r.FormValue("digest")

	b := blobs.GetBlob(name, digest)
	if b == nil {
		fmt.Printf("get blob by %s:%s error\n", name, digest)
		RenderErrorMsg(w, r, "blob of "+digest+" not exist")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", b.Size))
	n, err := w.Write(b.Data)
	if err != nil {
		fmt.Printf("Write to client error\n")
		return
	}

	fmt.Printf("Write to Client %d bytes blob OK\n", n)
}

// DeleteBlob to delete blob from api
func deleteBlob(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	digest := r.FormValue("digest")

	b := blobs.GetBlobPartial(name, digest)
	if b == nil {
		w.WriteHeader(http.StatusNotFound)
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
