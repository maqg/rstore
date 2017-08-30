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

	// RouteNameBlobUploadChunk for blobs upload management
	RouteNameBlobUploadChunk = "blob-upload-chunk"

	// RouteNameBlobsManifest for blobs and manifest relationship
	RouteNameBlobsManifest = "blobs-manifest"
)

const (
	// APIURLFormatBlobUpload url format for blobs upload, address:name
	APIURLFormatBlobUpload = "http://%s/v1/%s/blobs/uploads/"

	// APIURLFormatBlobsManifest for blobs-manifest operation address:name:blobsum
	APIURLFormatBlobsManifest = "http://%s/v1/%s/blobsmanifest/%s/"

	// APIURLFormatManifests for manifests management
	APIURLFormatManifests = "http://%s/v1/%s/manifests/%s/"
)

// EndPoints for all endpoints
var EndPoints = []string{
	RouteNameBase,
	RouteNameHelp,
	RouteNameHelpModule,
	RouteNameManifest,
	RouteNameBlobsManifest,
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
