package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/syariatifaris/goreportcard/cmd/goreportcard-cli/exec"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
	th      = flag.Float64("t", 0, "Threshold of failure command")
	format  = flag.Bool("f", false, "Print Using format (JSON)")
	hook    = flag.String("hook", "", "hook config file location")
	timeout = flag.Int("timeout", 5, "hook request timeout")
)

func main() {
	flag.Parse()
	e, err := exec.New(&exec.Config{
		Dir:         *dir,
		UseFormat:   *format,
		Verbose:     *verbose,
		FailThres:   *th,
		HookFile:    *hook,
		HookTimeout: *timeout,
	})
	if err != nil {
		fmt.Println("unable to execute:", err.Error())
		os.Exit(1)
	}
	if e == nil {
		log.Fatalln("fatal, unable to start checker. Invalid runner type")
	}
	ok, err := e.Run(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	if !ok {
		os.Exit(1)
		return
	}
	return
}
