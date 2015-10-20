package main

import (
	"flag"
	"./manager/configuration"
	"./manager/data"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"os"
	"html/template"
	"log"
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
	if _, err := os.Stat("dalek/manifest.json"); os.IsNotExist(err) {
		config.DebugLog("Makeing manifest.json")
		//asign default values to config.manifest
		json, err := json.MarshalIndent(config.Manifest, "", "  ")
		if(err != nil) {panic(err)}
		ioutil.WriteFile("dalek/manifest.json", json, 0775)
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

	data, err := ioutil.ReadFile("dalek/manifest.json")
	config.DebugErrorLog(err)
	err = json.Unmarshal(data, &config.Manifest)
	fmt.Println(config)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/", serveTemplate)
	http.ListenAndServe(":8080", nil)
}

func serveTemplate(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println(request.URL.Path)
	includesPath := path.Join("web", "dynamic", "includes.html")
	filePath := path.Join("web", "dynamic", request.URL.Path)
	data := data.PageWrapper{}
	if (request.URL.Path == "/") {
		filePath = path.Join("web", "dynamic", "index.html")
	}

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
		log.Println(err)
		http.Error(*writer, http.StatusText(code), code)
		return true
	}
	return false
}