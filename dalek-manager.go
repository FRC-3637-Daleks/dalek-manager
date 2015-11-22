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
	"strings"
	"bytes"
	"github.com/gorilla/context"
)

var config configuration.Config
var debug = flag.Bool("debug", false, "If set debug output will print")

func main() {
	flag.Parse()
	config.Debug = *debug
	if _, err := os.Stat("dalek"); os.IsNotExist(err) {
		config.DebugLog("Makeing dalek Directory")
		err := os.MkdirAll("dalek", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/autonomous"); os.IsNotExist(err) {
		config.DebugLog("Makeing autonomous Directory")
		err := os.MkdirAll("dalek/autonomous", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/ports"); os.IsNotExist(err) {
		config.DebugLog("Makeing ports Directory")
		err := os.MkdirAll("dalek/ports", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/controls"); os.IsNotExist(err) {
		config.DebugLog("Makeing controls Directory")
		err := os.MkdirAll("dalek/controls", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/settings"); os.IsNotExist(err) {
		config.DebugLog("Makeing settings Directory")
		err := os.MkdirAll("dalek/settings", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/logs"); os.IsNotExist(err) {
		config.DebugLog("Makeing logs Directory")
		err := os.MkdirAll("dalek/logs", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/binaries"); os.IsNotExist(err) {
		config.DebugLog("Makeing binaries Directory")
		err := os.MkdirAll("dalek/binaries", 0775)
		if (err != nil) {panic(err)}
	}
	if _, err := os.Stat("dalek/manifest.json"); os.IsNotExist(err) {
		config.DebugLog("Makeing manifest.json")
		json, err := json.MarshalIndent(config.Manifest, "", "  ")
		if (err != nil) {panic(err)}
		ioutil.WriteFile("dalek/manifest.json", json, 0664)
	} else {
		data, err := ioutil.ReadFile("dalek/manifest.json")
		config.DebugErrorLog(err)
		err = json.Unmarshal(data, &config.Manifest)
	}
	config.DebugLog("Loaded manifest: ", config.Manifest)
	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	rtr.HandleFunc("/", rootHandler)
	rtr.HandleFunc("/autonomous", autonomousHandler)
	rtr.HandleFunc("/ports", portsHandler)
	rtr.HandleFunc("/controls", controlsHandler)
	rtr.HandleFunc("/settings", settingsHandler)
	rtr.HandleFunc("/logs", logsHandler)
	rtr.HandleFunc("/binaries", binariesHandler)
	rtr.HandleFunc("/editor/{fileName}", editorHandler)
	rtr.HandleFunc("/editor/{fileType:autonomous|control|ports|settings}/{fileName}", editorHandler)
	rtr.HandleFunc("/file/{fileName}", fileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileName}", addFileHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileType:autonomous|control|ports|settings}/{fileName}", fileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileType:autonomous|control|ports|settings}/{fileName}", addFileHandler).Methods("POST")
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "index.html"), nil)
}

func autonomousHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "autonomous.html"), nil)
}

func portsHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "ports.html"), nil)
}

func controlsHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "controls.html"), nil)
}

func settingsHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "settings.html"), nil)
}

func logsHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "logs.html"), nil)
}

func binariesHandler(writer http.ResponseWriter, request *http.Request) {
	serveTemplate(writer, request, path.Join("web", "dynamic", "binaries.html"), nil)
}

func editorHandler(writer http.ResponseWriter, request *http.Request) {
	editorWrapper := data.EditorWrapper{}
	editorWrapper.Lang = "text"
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	filePath := "dalek/" + fileType + "/" + fileName
	config.DebugLog("Request for: ", filePath)
	editorWrapper.FileName = fileName
	editorWrapper.FileType = fileType
	fileList, err := getFileList("dalek/" + fileType)
	if(check(err, 500, &writer)) {return}
	fileJson, err := json.Marshal(fileList)
	if(check(err, 500, &writer)) {return}
	editorWrapper.Files = string(fileJson)
	config.DebugLog(editorWrapper.Files)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		config.DebugLog("Loading file into editor: ", filePath)
		content, err := ioutil.ReadFile(filePath)
		if (check(err, 500, &writer)) {return}
		editorWrapper.FileContent = string(content)
	}
	temp := strings.Split(fileName, ".")
	fileExt := temp[len(temp) - 1]
	editorWrapper.Lang = fileExt
	serveTemplate(writer, request, path.Join("web", "dynamic", "editor.html"), editorWrapper)
}

func addFileHandler(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if(check(err, 500, &writer)) {return }
	err = request.ParseMultipartForm(32 << 20)
	if(check(err, 500, &writer)) {return }
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	file, _, err := request.FormFile("file")
	if(check(err, 500, &writer)) {return }
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if(check(err, 500, &writer)) {return }
	ioutil.WriteFile("dalek/" + fileType + "/" + fileName, buf.Bytes(), 0664)
	config.DebugLog("Wrote file: " + fileType + " " + fileName)
	http.Error(writer, http.StatusText(200), 200)
}

func fileHandler(writer http.ResponseWriter, request *http.Request)  {
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	filePath := "dalek/" + fileType + "/" + fileName
	config.DebugLog("Request for: ", filePath)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		config.DebugLog("Serving file: ", filePath)
		content, err := ioutil.ReadFile(filePath)
		if (check(err, 500, &writer)) {return}
		writer.Write(content)
		return
	} else {
		config.DebugLog("File does not exist: ", filePath)
		http.Error(writer, http.StatusText(400), 400)
		return
	}
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
	if (check(err, 500, &writer)) {return }
	if err := tmpl.ExecuteTemplate(writer, "main", data); err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(500), 500)
	}
}

func check(err error, code int, writer *http.ResponseWriter) bool {
	if (err != nil) {
		config.DebugErrorLog(err)
		http.Error(*writer, http.StatusText(code), code)
		return true
	}
	return false
}

func getFileList(dir string) ([]string, error) {
	var (
		files []string
		err error
	)
	fileData, err := ioutil.ReadDir(dir)
	for _, f := range fileData {
		files = append(files, f.Name())
	}
	return files, err
}