package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"path"
	"os"
	"html/template"
	"log"

	"github.com/FRC-3637-Daleks/dalek-manager/manager/configuration"
	"github.com/FRC-3637-Daleks/dalek-manager/manager/data"

	"github.com/gorilla/mux"
)

var config configuration.Config
var debug  = flag.Bool("debug", false, "If set debug output will print")

func main() {
	flag.Parse()
	config.Debug = *debug
	if _, err := os.Stat("dalek"); os.IsNotExist(err) {
		config.DebugLog("Makeing dalek Directory")
		err := os.MkdirAll("dalek", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/autonomous"); os.IsNotExist(err) {
		config.DebugLog("Makeing autonomous Directory")
		err := os.MkdirAll("dalek/autonomous", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/ports"); os.IsNotExist(err) {
		config.DebugLog("Makeing ports Directory")
		err := os.MkdirAll("dalek/ports", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/controls"); os.IsNotExist(err) {
		config.DebugLog("Makeing controls Directory")
		err := os.MkdirAll("dalek/controls", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/settings"); os.IsNotExist(err) {
		config.DebugLog("Makeing settings Directory")
		err := os.MkdirAll("dalek/settings", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/logs"); os.IsNotExist(err) {
		config.DebugLog("Makeing logs Directory")
		err := os.MkdirAll("dalek/logs", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/binaries"); os.IsNotExist(err) {
		config.DebugLog("Makeing binaries Directory")
		err := os.MkdirAll("dalek/binaries", 0775)
		if(err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/manifest.json"); os.IsNotExist(err) {
		config.DebugLog("Makeing manifest.json")
		//asign default values to config.manifest
		json, err := json.MarshalIndent(config.Manifest, "", "  ")
		if(err != nil) {panic(err)}
		ioutil.WriteFile("dalek/manifest.json", json, 0775)
	} else {
		data, err := ioutil.ReadFile("dalek/manifest.json")
		config.DebugErrorLog(err)
		err = json.Unmarshal(data, &config.Manifest)
	}
	config.DebugLog("Loaded manifest: ", config.Manifest)
	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	rtr.HandleFunc("/", serveRoot)
	rtr.HandleFunc("/editor/{fileType:autonomous|ports|settings}/{fileName}", serveEditor).Methods("GET")
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}

func serveRoot(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "index.html"), data.PageWrapper{})
}

func serveEditor(writer http.ResponseWriter, request *http.Request)  {
	editorWrapper := data.EditorWrapper{}
	editorWrapper.Lang = "json"
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	config.DebugLog("Loading file into editor: ", "dalek/" + fileType + "/" + fileName)
	content, err := ioutil.ReadFile("dalek/" + fileType + "/" + fileName)
	if(check(err, 500, &writer)){return}
	editorWrapper.FileContent = string(content)
	if(fileType == "autonomous") {
		editorWrapper.Lang = "lua"
	}
	serveTemplate(writer, request, path.Join("web", "dynamic", "editor.html"), editorWrapper)
}

func serveTemplate(writer http.ResponseWriter, request *http.Request, filePath string, data interface{}) {
	includesPath := path.Join("web", "dynamic", "includes.html")
	config.DebugLog("Request for: \"", request.URL.Path, "\"\tSending: \"", filePath, "\"")
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
			return
		}
	}
	if info.IsDir() {
		http.NotFound(writer, request)
		return
	}
	tmpl, err := template.ParseFiles(includesPath, filePath)
	if(check(err, 500, &writer)){return }
	if err := tmpl.ExecuteTemplate(writer, "main", data); err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(500), 500)
	}
}

func check(err error, code int,  writer *http.ResponseWriter) bool {
	if(err != nil) {
		config.DebugErrorLog(err)
		http.Error(*writer, http.StatusText(code), code)
		return true
	}
	return false
}