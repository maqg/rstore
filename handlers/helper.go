package handlers

import (
	"fmt"
	"net/http"
	"octlink/rstore/api/v1"
	"octlink/rstore/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func apiHelp(w http.ResponseWriter, r *http.Request) {

	var data string
	module := mux.Vars(r)["module"]
	if module == "" {
		data = utils.JSON2String(v1.EndPoints)
	} else {
		data = utils.JSON2String(v1.RouteDescriptorsMap[module])
	}

	fmt.Fprint(w, data)
}

// appropriate handler for handling image manifest requests.
func helperManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET": http.HandlerFunc(apiHelp),
	}

	return mhandler
}
