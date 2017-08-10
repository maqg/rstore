package handlers

import (
	"net/http"
	"octlink/rstore/modules/blobs"

	"github.com/gorilla/handlers"
)

func blobUploadManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(blobs.GetUploadStatus),
		"HEAD": http.HandlerFunc(blobs.GetUploadStatus),
	}

	mhandler["POST"] = http.HandlerFunc(blobs.StartBlobUpload)
	mhandler["PATCH"] = http.HandlerFunc(blobs.PatchBlobData)
	mhandler["PUT"] = http.HandlerFunc(blobs.PutBlobUploadComplete)
	mhandler["DELETE"] = http.HandlerFunc(blobs.CancelBlobUpload)

	return mhandler
}
