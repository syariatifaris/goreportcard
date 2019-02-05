package main

import (
	"flag"
	"fmt"
	"os"

	//"fmt"
	"log"

	"github.com/syariatifaris/goreportcard/cmd/goreportcard-cli/exec"
	//"os"
	//"github.com/gojp/goreportcard/check"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
	th      = flag.Float64("t", 0, "Threshold of failure command")
	format  = flag.Bool("f", false, "Print Using format (JSON)")
)

func main() {
	flag.Parse()
	e := exec.New(&exec.Config{Dir: *dir, UseFormat: *format, Verbose: *verbose, FailThres: *th})
	if e == nil {
		log.Fatalln("Fatal, unable to start checker. Invalid runner type")
	}
	err := e.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return
}
