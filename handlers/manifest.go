package handlers

import (
	"net/http"
	"octlink/rstore/modules/manifest"

	"github.com/gorilla/handlers"
)

func manifestManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":    http.HandlerFunc(manifest.GetManifest),
		"PUT":    http.HandlerFunc(manifest.PutManifest),
		"DELETE": http.HandlerFunc(manifest.DeleteManifest),
	}

	return mhandler
}
