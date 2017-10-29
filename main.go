package main

import (
	"net/http"
)

func main() {
	var router = Route{}

	router.addRoute("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Root Page"))
	})

	http.ListenAndServe(":8000", router.serveRequests())
}
