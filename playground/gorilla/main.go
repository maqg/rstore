package main

import (
	"log"
	"net/http"
	"octlink/rstore/playground/gorilla/v1"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range v1.RouteDescriptorsMap {
		w.Write([]byte(k))
		w.Write([]byte(v.Path))
	}

	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r := v1.NewRouters()

	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
