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
	data := data.DataWrapper{}
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