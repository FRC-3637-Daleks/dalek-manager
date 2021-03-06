/* Team 3637 Dalek Manager - A web base robot configuration manager
   Copyright (C) 2016  Team 3637

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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
	"github.com/FRC-3637-Daleks/dalek-manager/manager/util"
	"errors"
)

var config configuration.Config
var manifest model.Manifest
var debug = flag.Bool("debug", false, "If set debug output will print")

func main() {
	flag.Parse()
	config.Debug = *debug
	if _, err := os.Stat("dalek/manifest.json"); os.IsNotExist(err) {
		config.DebugLog("Makeing manifest.json")
		manifest.Server.Port = 5810
		manifest.Server.MqttPort = 5800
		manifest.Server.WebRoot = "web/"
		manifest.Templates.Configs.Controls = "controls/schema.json"
		manifest.Templates.Configs.Ports = "ports/schema.json"
		manifest.Templates.Configs.Settings = "settings/schema.json"
		json, err := json.MarshalIndent(manifest, "", "    ")
		if (err != nil) {panic(err)}
		ioutil.WriteFile("dalek/manifest.json", json, 0664)
	} else {
		data, err := ioutil.ReadFile("dalek/manifest.json")
		config.ErrorLog(err)
		err = json.Unmarshal(data, &manifest)
	}
	json, err := json.MarshalIndent(manifest, "", "    ")
	if(err != nil) {
		config.Log("Manifest error: ", manifest)
		config.ErrorLog(err)
		return
	}
	config.DebugLog("Loaded manifest: \n", string(json))
	if _, err := os.Stat(manifest.Server.WebRoot); os.IsNotExist(err) {
		panic(errors.New("Cannot find web directory"))
	}
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
	fileRegex := "autonomous|controls|ports|settings|logs|binaries"
	editorRegex := "autonomous|controls|ports|settings"
	editorGuiRegex := "controls|ports|settings"
	rtr := mux.NewRouter()
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir(path.Join(manifest.Server.WebRoot, "static")))))
	rtr.HandleFunc("/", rootHandler)
	rtr.HandleFunc("/{fileType:" + fileRegex + "}", defaultHandler)
	rtr.HandleFunc("/binaries/pull", pullBinHandler)
	rtr.HandleFunc("/editor-gui/{fileType:" + editorGuiRegex + "}/{fileName}", configHandler)
	rtr.HandleFunc("/log-viewer/{fileName}", logViewerHandler)
	rtr.HandleFunc("/editor/{fileName}", editorHandler).Methods("GET")
	rtr.HandleFunc("/editor/{fileName}", putEditorHandler).Methods("POST")
	rtr.HandleFunc("/editor/{fileType:" + editorRegex + "}/{fileName}", editorHandler).Methods("GET")
	rtr.HandleFunc("/editor/{fileType:" + editorRegex + "}/{fileName}", putEditorHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileName}", getFileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileName}", addFileHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileName}", deleteFileHandler).Methods("DELETE")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", getFileHandler).Methods("GET")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", addFileHandler).Methods("POST")
	rtr.HandleFunc("/file/{fileType:" + fileRegex + "}/{fileName}", deleteFileHandler).Methods("DELETE")
	rtr.HandleFunc("/file/list/", listFileHandler).Methods("GET")
	rtr.HandleFunc("/file/list/{fileType:" + fileRegex + "}", listFileHandler).Methods("GET")
	rtr.HandleFunc("/status", statusHandler)
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
		config.Log("Starting server on port: " + strconv.Itoa(manifest.Server.Port))
		http.ListenAndServe(":" + strconv.Itoa(manifest.Server.Port), context.ClearHandler(http.DefaultServeMux))
	}
	defaultStart:
	config.Log("Starting server on port: 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	var pageData data.PageWrapper
	ip, err := getMqttHost()
	if (check(err, 500, &writer)) {return}
	pageData.ServerIP = ip
	pageData.MqttPort = manifest.Server.MqttPort
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", "index.html"), pageData)
}

func defaultHandler(writer http.ResponseWriter, request *http.Request) {
	var pageData data.PageWrapper
	ip, err := getMqttHost()
	if (check(err, 500, &writer)) {return}
	pageData.ServerIP = ip
	pageData.MqttPort = manifest.Server.MqttPort
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", fileType + ".html"), pageData)
}

func configHandler(writer http.ResponseWriter, request *http.Request)  {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", "config.html"), nil)
}

func logViewerHandler(writer http.ResponseWriter, request *http.Request)  {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", "log-viewer.html"), nil)
}

func editorHandler(writer http.ResponseWriter, request *http.Request) {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	var editorWrapper data.EditorWrapper
	ip, err := getMqttHost()
	if (check(err, 500, &writer)) {return}
	editorWrapper.ServerIP = ip
	editorWrapper.MqttPort = manifest.Server.MqttPort
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
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", "editor.html"), editorWrapper)
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
	serveTemplate(writer, request, path.Join(manifest.Server.WebRoot, "dynamic", "editor.html"), editorWrapper)
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
	vars := mux.Vars(request)
	fileType := vars["fileType"]
	var fileList []string
	var err error
	if(fileType == "logs") {
		fileList, err = getLogList()
	} else {
		fileList, err = getFileList("dalek/" + fileType)
	}
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
	if (fileType == "" && fileName == "manifest.json") {
		data, err := ioutil.ReadFile("dalek/manifest.json")
		config.ErrorLog(err)
		err = json.Unmarshal(data, &manifest)
		if (check(err, 500, &writer)) {return }
		if(manifest.Runtime.Binary != "") {
			config.Log("Setting binary: " + manifest.Runtime.Binary)
			if (check(err, 500, &writer)) {return }
			err = util.CopyFile("dalek/" + manifest.Runtime.Binary, "../FRCUserProgram")
			if (check(err, 500, &writer)) {return }
		}
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func deleteFileHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	file := "dalek/" + vars["fileType"] + "/" + vars["fileName"]
	config.DebugLog("Deleting file: " + file)
	err := os.Remove(file)
	if (check(err, 500, &writer)) {return }
	writer.WriteHeader(http.StatusNoContent)
	writer.Write([]byte("No Content"))
}

func pullBinHandler(writer http.ResponseWriter, request *http.Request) {
	config.DebugLog("Request for: " + request.Method + " \"", request.URL.Path, "\"")
	err := request.ParseForm()
	if (check(err, 500, &writer)) {return }
	config.DebugLog(request.Form)
	config.DebugLog(request.PostForm)
	fileName := request.PostFormValue("fileName")
	if(fileName == "") {
		err = errors.New("File name cannot be empty")
		if (check(err, 500, &writer)) {return }
	}
	config.DebugLog("Copying FRCUserProgram to : dalek/binaries/" + fileName)
	err = util.CopyFile("../FRCUserProgram", "dalek/binaries/" + fileName)
	if (check(err, 500, &writer)) {return }
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func statusHandler(writer http.ResponseWriter, request *http.Request)  {
	writer.WriteHeader(http.StatusOK)
}

func serveTemplate(writer http.ResponseWriter, request *http.Request, filePath string, data interface{}) {
	includesPath := path.Join(manifest.Server.WebRoot, "dynamic", "includes.html")
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
		config.ErrorLog(err)
		http.Error(*writer, err.Error(), code)
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

func getLogList() ([]string, error) {
	var (
		files []string
		err error
	)
	fileData, err := ioutil.ReadDir("dalek/logs")
	for _, f := range fileData {
		if (f.IsDir()) {
			console, manifest := false, false
			fileData, err := ioutil.ReadDir("dalek/logs/" + f.Name())
			if (err != nil) {
				return files, err
			}
			for _, f := range fileData {
				if(f.Name() == "manifest.json") {
					manifest = true
				} else if (f.Name() == "console.log") {
					console = true
				}
			}
			if (manifest && console) {
				files = append(files, f.Name())
			}
		}
	}
	return files, err
}

func getMqttHost() (string, error) {
	if(manifest.Server.MqttHost != "") {
		return manifest.Server.MqttHost, nil
	} else {
		return util.GetIP()
	}
}