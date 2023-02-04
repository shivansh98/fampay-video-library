package server

import (
	"net/http"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-videos", getVideoMiddleWare)
	mux.HandleFunc("/search", SearchmiddleWare)
	http.ListenAndServe(":8080", mux)
}
