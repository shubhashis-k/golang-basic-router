package main

import (
	"net/http"
)
type WebHandler = func(w http.ResponseWriter, r *http.Request)
type MiddleWareFunction = func(w http.ResponseWriter, r *http.Request) bool
type WebNode struct{
	method string
	pattern string
	staticFileLocation string
	//Dynamic params?
}

type Route struct{
	routeMap map[WebNode] WebHandler
	middlewareFunctions []MiddleWareFunction
}

func (route *Route) addRoute(method string, pattern string, handler func(w http.ResponseWriter, r *http.Request)){
	//No additional checks yet
	if(route.routeMap == nil){
		route.routeMap = make(map[WebNode] WebHandler)
	}
	var webnode = WebNode{method: method, pattern: pattern}
	route.routeMap[webnode] = handler
}

func (route *Route) handleMiddleWareFunction(w http.ResponseWriter, r *http.Request) bool{
	for i:= 0 ; i < len(route.middlewareFunctions) ; i++{
		var next = route.middlewareFunctions[i](w, r)

		if(!next){
			return false
		}
	}
	return true
}

func (route *Route) addMiddleWareFunction(middleWareFunction func(w http.ResponseWriter, r *http.Request) bool){
	route.middlewareFunctions = append(route.middlewareFunctions, middleWareFunction)
}


func (route *Route) serveRequests() http.HandlerFunc{
	var handlerFunction = func (w http.ResponseWriter, r * http.Request){
		//go through middleware functions
		if(!route.handleMiddleWareFunction(w, r)){
			return
		}
		//go through route functions
		var webnode = WebNode{method:r.Method, pattern:r.URL.Path}
		if(route.routeMap[webnode] != nil){
			route.routeMap[webnode](w, r)
			return
		} else{
			w.Write([]byte("404"))
		}
	}

	var handler = http.HandlerFunc(handlerFunction)

	return handler
}