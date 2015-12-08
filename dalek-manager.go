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
	"bytes"
	"github.com/gorilla/context"
	"github.com/FRC-3637-Daleks/dalek-manager/manager/model"
	"strconv"
)

var config configuration.Config
var manifest model.Manifest
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
		manifest.Server.Port = 8080
		json, err := json.MarshalIndent(manifest, "", "    ")
		if (err != nil) {panic(err)}
		ioutil.WriteFile("dalek/manifest.json", json, 0664)
	} else {
		data, err := ioutil.ReadFile("dalek/manifest.json")
		config.DebugErrorLog(err)
		err = json.Unmarshal(data, &manifest)
	}
	fileRegex := "autonomous|controls|ports|settings|logs|binaries"
	config.DebugLog("Loaded manifest: ", manifest)
	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	rtr.HandleFunc("/", rootHandler)
	rtr.HandleFunc("/autonomous", autonomousHandler)
	rtr.HandleFunc("/ports", portsHandler)
	rtr.HandleFunc("/controls", controlsHandler)
	rtr.HandleFunc("/settings", settingsHandler)
	rtr.HandleFunc("/logs", logsHandler)
	rtr.HandleFunc("/binaries", binariesHandler)
	rtr.HandleFunc("/editor/{fileName}", editorHandler).Methods("GET")
	rtr.HandleFunc("/editor/{fileName}", putEditorHandler).Methods("POST")
	rtr.HandleFunc("/editor/{fileType:" + fileRegex + "}/{fileName}", editorHandler).Methods("GET")
	rtr.HandleFunc("/editor/{fileType:" + fileRegex + "}/{fileName}", putEditorHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileName}", getFileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileName}", addFileHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileName}", deleteFileHandler).Methods("DELETE")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", getFileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", addFileHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", deleteFileHandler).Methods("DELETE")
	rtr.HandleFunc("/file/list/", listFileHandler).Methods("GET")
	rtr.HandleFunc("/file/list/{fileType:" + fileRegex + "}", listFileHandler).Methods("GET")
	http.Handle("/", rtr)
	if (manifest.Server.Port == 0) {
		goto defaultStart
	} else {
		if (manifest.Server.Port < 1) {
			config.DebugLog("Port: " + strconv.Itoa(manifest.Server.Port) + "Is not a valid port")
			goto defaultStart
		}
		if (manifest.Server.Port < 1024) {
			config.DebugLog("Please use a port higher than 1023")
			goto defaultStart
		}
		config.DebugLog("Starting server on port: " + strconv.Itoa(manifest.Server.Port))
		http.ListenAndServe(":" + strconv.Itoa(manifest.Server.Port), context.ClearHandler(http.DefaultServeMux))
	}
	defaultStart:
	config.DebugLog("Starting server on port: 8080")
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
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	editorWrapper := data.EditorWrapper{}
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	filePath := "dalek/" + fileType + "/" + fileName
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		config.DebugLog("Loading file into editor: ", filePath)
		content, err := ioutil.ReadFile(filePath)
		if (check(err, 500, &writer)) {return}
		editorWrapper.FileContent = string(content)
	}
	serveTemplate(writer, request, path.Join("web", "dynamic", "editor.html"), editorWrapper)
}

func putEditorHandler(writer http.ResponseWriter, request *http.Request) {
	editorWrapper := data.EditorWrapper{}
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	filePath := "dalek/" + fileType + "/" + fileName
	err := request.ParseForm()
	if (check(err, 500, &writer)) {return }
	err = request.ParseMultipartForm(32 << 20)
	if (check(err, 500, &writer)) {return }
	config.DebugLog("Request for: ", filePath)
	file, _, err := request.FormFile("file")
	if (check(err, 500, &writer)) {return }
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if (check(err, 500, &writer)) {return }
	editorWrapper.FileContent = string(buf.Bytes())
	serveTemplate(writer, request, path.Join("web", "dynamic", "editor.html"), editorWrapper)
}

func getFileHandler(writer http.ResponseWriter, request *http.Request) {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	filePath := "dalek/" + fileType + "/" + fileName
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

func listFileHandler(writer http.ResponseWriter, request *http.Request) {
	fileList, err := getFileList("dalek/" + request.RequestURI[11:len(request.RequestURI)])
	if (check(err, 500, &writer)) {return}
	fileJson, err := json.MarshalIndent(fileList, "", "    ")
	if (check(err, 500, &writer)) {return}
	writer.Write(fileJson)
}

func addFileHandler(writer http.ResponseWriter, request *http.Request) {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	err := request.ParseForm()
	if (check(err, 500, &writer)) {return }
	err = request.ParseMultipartForm(32 << 20)
	if (check(err, 500, &writer)) {return }
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	fileName := vars["fileName"]
	file, _, err := request.FormFile("file")
	if (check(err, 500, &writer)) {return }
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if (check(err, 500, &writer)) {return }
	ioutil.WriteFile("dalek/" + fileType + "/" + fileName, buf.Bytes(), 0664)
	config.DebugLog("Wrote file: " + fileType + " " + fileName)
	http.Error(writer, http.StatusText(200), 200)
}

func deleteFileHandler(writer http.ResponseWriter, request *http.Request) {
	file := "dalek/" + request.RequestURI[6:len(request.RequestURI)]
	config.DebugLog("Deleting file: " + file)
	err := os.Remove(file)
	if (check(err, 500, &writer)) {return }
	writer.WriteHeader(http.StatusNoContent)
	writer.Write([]byte("No Content"))
}

func serveTemplate(writer http.ResponseWriter, request *http.Request, filePath string, data interface{}) {
	includesPath := path.Join("web", "dynamic", "includes.html")
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	config.DebugLog("Sending: \"", filePath, "\"")
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
		if (!f.IsDir()) {
			files = append(files, f.Name())
		}
	}
	return files, err
}