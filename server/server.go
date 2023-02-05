package server

import (
	"net/http"
)

// Run runs the mux server and initializes the routes with the handlerFuncs
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-videos", getVideoMiddleWare)
	mux.HandleFunc("/search", SearchMiddleWare)
	http.ListenAndServe(":8080", mux)
}
