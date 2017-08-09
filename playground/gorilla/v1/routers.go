package v1

import "github.com/gorilla/mux"

func NewRouters() *mux.Router {

	router := mux.NewRouter()

	router.StrictSlash(true)

	for _, descriptor := range routeDescriptors {
		router.Path(descriptor.Path).Name(descriptor.Name)
	}

	return router
}
