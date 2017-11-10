package CustomFileServer

import (
	"net/http"
	"os"
	"log"
	"fmt"
)

type WebHandler = func(w http.ResponseWriter, r *http.Request)
type MiddleWareFunction = func(w http.ResponseWriter, r *http.Request) bool
type FileNode struct{
	fakeDirectory string
	actualDirectory string
	rootFile string
	//Dynamic params?
}

type FileRoute struct{
	fileMap map[string] string
	rootMap map[string] string
	middlewareFunctions []MiddleWareFunction
}

func (fileRoute *FileRoute) AddFileNode(fakeDirectory string, actualDirectory string, rootFile string){
	//No additional checks yet
	if(fileRoute.fileMap == nil){
		fileRoute.fileMap = make(map[string] string)
	}
	if(fileRoute.rootMap == nil){
		fileRoute.rootMap = make(map[string] string)
	}
	fileRoute.fileMap[fakeDirectory] = actualDirectory
	fileRoute.rootMap[fakeDirectory] = rootFile
}

func (fileRoute *FileRoute) handleMiddleWareFunction(w http.ResponseWriter, r *http.Request) bool{
	for i:= 0 ; i < len(fileRoute.middlewareFunctions) ; i++{
		var next = fileRoute.middlewareFunctions[i](w, r)

		if(!next){
			return false
		}
	}
	return true
}

func (fileRoute *FileRoute) AddMiddleWareFunction(middleWareFunction func(w http.ResponseWriter, r *http.Request) bool){
	fileRoute.middlewareFunctions = append(fileRoute.middlewareFunctions, middleWareFunction)
}

func getDirectories(fullPath string) []string{
	s := fullPath + "/"
	var dirTree = []string{}
	var currDir = ""

	for i:= 0; i<len(s); i++{
		if(string(s[i]) == "/"){
			if(currDir != ""){
				dirTree = append(dirTree, currDir)
			}
			currDir = ""
		} else{
			currDir += string(s[i])
		}
	}

	return dirTree
}

func (fileRoute *FileRoute)resolveFileLocation(url string) string{
	//return the file location, after remapping. If there are no files, return an empty string

	directories := getDirectories(url)
	directories = append(directories, "" , "")
	//HOME or main page. Check for main page associations
	pageMainDirectory := directories[0]
	pageResourcesDirectory := ""

	for i := 1; i < len(directories) ; i++{
		pageResourcesDirectory += directories[i]
	}

	//WARNING. EMPTY STRING MIGHT BE PROBLEMATIC
	remappedDirectory := fileRoute.fileMap[pageMainDirectory]

	if(remappedDirectory == "") {
		return ""
	}

	if(pageResourcesDirectory == ""){
		return "./" + remappedDirectory + "/" + fileRoute.rootMap[pageMainDirectory]
	} else {
		return "./" + remappedDirectory + "/" + pageResourcesDirectory
	}
}
func (fileRoute *FileRoute) ServeRequests() http.HandlerFunc{
	var handlerFunction = func (w http.ResponseWriter, r * http.Request){
		//go through middleware functions
		if(!fileRoute.handleMiddleWareFunction(w, r)){
			return
		}
		//go through route functions
		//beware, only GET requests. have a check for that
		var mainFileLocation = fileRoute.resolveFileLocation(r.URL.Path)

		_,err := os.Stat(mainFileLocation)

		if(err == nil){
			http.ServeFile(w,r,mainFileLocation)
		} else{
			fmt.Println("Invalid File Location " + mainFileLocation)
		}

	}

	var handler = http.HandlerFunc(handlerFunction)

	return handler
}
