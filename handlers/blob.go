package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func getBlob(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func deleteBlob(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{\"msg\":\"this is blob message\"}"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))
	fmt.Fprint(w, emptyJSON)
}

func blobManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(getBlob),
		"HEAD": http.HandlerFunc(getBlob),
	}

	mhandler["DELETE"] = http.HandlerFunc(deleteBlob)

	return mhandler
}
