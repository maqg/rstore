package v1

import "github.com/gorilla/mux"

const (
	RouteNameBase            = "base"
	RouteNameHelp            = "help"
	RouteNameManifest        = "manifest"
	RouteNameTags            = "tags"
	RouteNameBlob            = "blob"
	RouteNameBlobUpload      = "blob-upload"
	RouteNameBlobUploadChunk = "blob-upload-chunk"
	RouteNameCatalog         = "catalog"
)

func NewRouters() *mux.Router {

	router := mux.NewRouter()

	router.StrictSlash(true)

	for _, descriptor := range routeDescriptors {
		router.Path(descriptor.Path).Name(descriptor.Name)
	}

	return router
}
