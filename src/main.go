package main

import (
	"flag"
	"fmt"
)

var generate = flag.Bool("gen", false, "If set will generate nonexistant configs")

func main() {
	flag.Parse()
	fmt.Println(*generate)
	//Generate Files

}