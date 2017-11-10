package main

import (
	//"net/http"
	//"fmt"
	//"strings"
	//"net/http"
	"net/http"
	"awesomeProject/CustomFileServer"
)

func main() {
	fileServer := CustomFileServer.FileRoute{}
	fileServer.AddFileNode("","StaticFiles","index.html")
	http.ListenAndServe(":8080",fileServer.ServeRequests())
}
