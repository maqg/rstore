package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/api/v1"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
)

func apiHelp(w http.ResponseWriter, r *http.Request) {
	data := utils.JSON2String(v1.RouteDescriptors)
	fmt.Fprint(w, data)
}

// appropriate handler for handling image manifest requests.
func helperManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET": http.HandlerFunc(apiHelp),
	}

	return mhandler
}
