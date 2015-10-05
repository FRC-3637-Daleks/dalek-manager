package main

import (
	"flag"
	"manager/configuration"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
}