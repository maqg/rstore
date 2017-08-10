package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

// apiBase implements a simple yes-man for doing overall checks against the
// api. This can support auth roundtrips to support docker login.
func apiHelp(w http.ResponseWriter, r *http.Request) {

	const emptyJSON = "{\"msg\":\"this is help message\"}"
	// Provide a simple /v2/ 200 OK response with empty json response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))

	fmt.Fprint(w, emptyJSON)
}

// appropriate handler for handling image manifest requests.
func helperManager(r *http.Request) http.Handler {

	mhandler := handlers.MethodHandler{
		"GET":  http.HandlerFunc(apiHelp),
		"HEAD": http.HandlerFunc(apiHelp),
	}

	mhandler["PUT"] = http.HandlerFunc(apiHelp)
	mhandler["DELETE"] = http.HandlerFunc(apiHelp)

	return mhandler
}
