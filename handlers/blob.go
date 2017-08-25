package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/modules/blobs"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func writeBackError(w http.ResponseWriter, r *http.Request) {
	errorMsg := "{\"status\":\"Error\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(errorMsg)))
	fmt.Fprint(w, errorMsg)
}

// GetBlob to get blob from web api
func getBlob(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]
	// digest := r.FormValue("digest")

	data, len, err := blobs.GetBlob(name, digest)
	if err != nil {
		fmt.Printf("get blob by %s:%s error\n", name, digest)
		writeBackError(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len))
	n, err := w.Write(data)
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

	var data string
	err := blobs.DeleteBlob(name, digest)
	if err != nil {
		data = "{\"status\":\"Error\"}"
	} else {
		data = "{\"status\":\"OK\"}"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, data)
}

func blobManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(getBlob),
		"HEAD": http.HandlerFunc(getBlob),
	}

	mhandler["DELETE"] = http.HandlerFunc(deleteBlob)

	return mhandler
}
