package main

import (
	"net/http"
	"io/ioutil"
	//"fmt"
)

func mainHandler(w http.ResponseWriter, r *http.Request){
	data, err := ioutil.ReadFile("./StaticFiles/index.html")

	if(err == nil){
		w.Write(data)
	}

}

func main() {
	var router = Route{}

	router.addfileServer("StaticFiles")
	router.addRoute("/", mainHandler)

	router.addRoute("/home", func(w http.ResponseWriter, r *http.Request){
		http.ServeFile(w, r, "/index.html")
	})

	http.ListenAndServe(":9090", router.serveRequests())
	//http.ListenAndServe(":8080", handler)
	//data, err := ioutil.ReadFile("index.html")
	//
	//if(err != nil){
	//	fmt.Println(data)
	//}
}
