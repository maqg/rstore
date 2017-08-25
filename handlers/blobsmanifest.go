package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getBlobsManifest(w http.ResponseWriter, r *http.Request) {

	blobsum := mux.Vars(r)["blobsum"]

	blobs := blobsmanifest.GetBlobsManifest(blobsum)
	if blobs == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := utils.JSON2String(blobs)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

func blobsmanifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET": http.HandlerFunc(getBlobsManifest),
	}

	return mhandler
}
