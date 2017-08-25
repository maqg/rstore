package v1

import "github.com/gorilla/mux"

const (
	// RouteNameBase base router
	RouteNameBase = "base"

	// RouteNameHelp helper router
	RouteNameHelp = "help"

	// RouteNameHelpModule helper router
	RouteNameHelpModule = "help-module"

	// RouteNameManifest manifest routers
	RouteNameManifest = "manifest"

	// RouteNameBlob for blob management
	RouteNameBlob = "blob"

	// RouteNameBlobUpload for blobs upload management
	RouteNameBlobUpload = "blob-upload"

	// RouteNameBlobUploadChunk for blob upload by chunk
	RouteNameBlobUploadChunk = "blob-upload-chunk"
)

// EndPoints for all endpoints
var EndPoints = []string{
	RouteNameBase,
	RouteNameHelp,
	RouteNameHelpModule,
	RouteNameManifest,
	RouteNameBlob,
	RouteNameBlobUpload,
	RouteNameBlobUploadChunk,
}

// NewRouters to new routes manager for http
func NewRouters() *mux.Router {

	router := mux.NewRouter()

	router.StrictSlash(true)

	for _, descriptor := range RouteDescriptors {
		router.Path(descriptor.path).Name(descriptor.Name)
	}

	return router
}
