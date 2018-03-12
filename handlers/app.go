package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"octlink/rstore/api/v1"
	"octlink/rstore/utils/octlog"

	"github.com/gorilla/mux"
)

var logger *octlog.LogConfig

// InitLog for handlers
func InitLog(level int) {
	logger = octlog.InitLogConfig("handlers.log", level)
}

// App is a global registry application object. Shared resources can be placed
// on this object that will be accessible from all requests. Any writable
// fields should be protected.
type App struct {
	Router *mux.Router // main application router, configured with dispatchers
	// httpHost is a parsed representation of the http.host parameter from
	// the configuration. Only the Scheme and Host fields are used.
	httpHost url.URL

	readOnly bool
}

// apiBase implements a simple yes-man for doing overall checks against the
// api. This can support auth roundtrips to support docker login.
func apiBase(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{}"
	// Provide a simple /v2/ 200 OK response with empty json response.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))

	fmt.Fprint(w, emptyJSON)
}

// for the route. The dispatcher will use this to dynamically create request
// specific handlers for each endpoint without creating a new router for each
// request.
type dispatchFunc func(r *http.Request) http.Handler

// dispatcher returns a handler that constructs a request specific context and
// handler, using the dispatch factory function.
func (app *App) dispatcher(dispatch dispatchFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dispatch(r).ServeHTTP(w, r)
	})
}

// register a handler with the application, by route name. The handler will be
// request time.
func (app *App) register(routeName string, dispatch dispatchFunc) {
	if app.Router == nil {
		logger.Errorf("app router is null\n")
	}
	app.Router.GetRoute(routeName).Handler(app.dispatcher(dispatch))
}

// NewApp to new an application
func NewApp() *App {

	app := &App{
		Router: v1.NewRouters(),
	}

	// Register the handler dispatchers.
	app.register(v1.RouteNameBase, func(r *http.Request) http.Handler {
		logger.Errorf("register base routers error\n")
		return http.HandlerFunc(apiBase)
	})

	app.register(v1.RouteNameHelp, helperManager)
	app.register(v1.RouteNameHelpModule, helperManager)
	app.register(v1.RouteNameManifest, manifestManager)
	app.register(v1.RouteNameBlob, blobManager)
	app.register(v1.RouteNameBlobUpload, blobUploadManager)
	app.register(v1.RouteNameBlobsManifest, blobsmanifestManager)

	logger.Infof("Created App and registered all routes OK\n")

	return app
}
