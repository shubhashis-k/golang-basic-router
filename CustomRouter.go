package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"strings"
)

type Route struct{
	routeMap map[string] func(w http.ResponseWriter, r *http.Request)
	fileServerLocation string
}

func (route *Route) addRoute(pattern string, handler func(w http.ResponseWriter, r *http.Request)){
	//No additional checks yet
	if(route.routeMap == nil){
		fmt.Println("Map is nil")
		route.routeMap = make(map[string] func(w http.ResponseWriter, r *http.Request))
	}
	route.routeMap[pattern] = handler
}

func (route *Route) addfileServer(directory string){
	route.fileServerLocation = directory
}


func (route *Route) serveRequests() http.HandlerFunc{
	fmt.Println("Called ServeRequest")



	var directoryToLookForFiles = "./" + route.fileServerLocation




	//might be problematic
	var handlerFunction = func (w http.ResponseWriter, r * http.Request){

		fmt.Println(r.RequestURI)

		//.js,.html
		if(strings.ContainsAny(r.RequestURI,".")){
			filedata, err := ioutil.ReadFile(directoryToLookForFiles + r.RequestURI)

			if(err == nil){
				w.Write(filedata)
			}
		}else if(route.routeMap[r.RequestURI] != nil) {
			route.routeMap[r.RequestURI](w, r)
		}else{
			fmt.Println("route not found")
		}
	}

	var handler = http.HandlerFunc(handlerFunction)

	return handler
}