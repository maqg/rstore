package handlers

import (
	"net/http"
	"octlink/rstore/modules/blobs"

	"github.com/gorilla/handlers"
)

func blobManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(blobs.GetBlob),
		"HEAD": http.HandlerFunc(blobs.GetBlob),
	}

	mhandler["DELETE"] = http.HandlerFunc(blobs.DeleteBlob)

	return mhandler
}
