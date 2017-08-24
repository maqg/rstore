package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octlink/rstore/modules/manifest"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// GetManifest for api call
func getManifest(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	digest := mux.Vars(r)["digest"]

	manifest := manifest.GetManifest(name, digest)
	if manifest == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(manifest)
	dataStr := utils.BytesToString(data)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(dataStr)))
	fmt.Fprint(w, dataStr)
}

// DeleteManifest for api call
func deleteManifest(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

// PutManifest for api call
func putManifest(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func manifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":    http.HandlerFunc(getManifest),
		"HEAD":   http.HandlerFunc(getManifest),
		"PUT":    http.HandlerFunc(putManifest),
		"DELETE": http.HandlerFunc(deleteManifest),
	}

	return mhandler
}
