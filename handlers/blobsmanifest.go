package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"octlink/rstore/modules/blobsmanifest"
	"octlink/rstore/utils"
	"octlink/rstore/utils/octlog"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getBlobsManifest(w http.ResponseWriter, r *http.Request) {

	blobsum := mux.Vars(r)["digest"]
	if blobsum == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

func postBlobsManifest(w http.ResponseWriter, r *http.Request) {

	blobsum := mux.Vars(r)["digest"]
	size := r.FormValue("size")

	if size == "" || blobsum == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	blobs := blobsmanifest.GetBlobsManifest(blobsum)
	if blobs != nil {
		octlog.Warn("blobsmanifest of %s already exist\n", blobsum)
		w.WriteHeader(http.StatusOK)
		return
	}

	bm := new(blobsmanifest.BlobsManifest)
	bm.BlobSum = blobsum

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		octlog.Error("read all data from r.Body error\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.Unmarshal(data, &bm.Chunks)
	bm.Size = utils.StringToInt64(size)

	err = bm.Write()
	if err != nil {
		octlog.Error("Write blobs-manifest to server error\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	octlog.Debug("write new blobs-manifest %s OK\n", bm.BlobSum)
}

func blobsmanifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(getBlobsManifest),
		"POST": http.HandlerFunc(postBlobsManifest),
	}

	return mhandler
}
