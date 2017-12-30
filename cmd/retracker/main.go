package main

import (
	Core "github.com/vvampirius/retracker/core"
	"github.com/vvampirius/retracker/core/common"
	"flag"
	"fmt"
	"syscall"
	"os"
)

const VERSION  = 0.3

func PrintRepo(){
	fmt.Fprintln(os.Stderr, "\n# https://github.com/vvampirius/retracker")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		PrintRepo()
	}
	listen := flag.String("l", ":80", "Listen address:port")
	age := flag.Float64("a", 180, "Keep 'n' minutes peer in memory")
	debug := flag.Bool("d", false, "Debug mode")
	xrealip := flag.Bool("x", false, "Get RemoteAddr from X-Real-IP header")
	forwards := flag.String("f", "", "Load forwards from YAML file")
	forwardTimeout := flag.Int("t", 2, "Timeout (sec) for forward requests (used with -f)")
	ver := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *ver {
		fmt.Println(VERSION)
		PrintRepo()
		syscall.Exit(0)
	}

	config := common.Config{
		Listen: *listen,
		Debug: *debug,
		Age: *age,
		XRealIP: *xrealip,
		ForwardTimeout: *forwardTimeout,
	}

	if *forwards != `` {
		if err := config.ReloadForwards(*forwards); err!=nil {
			fmt.Fprintln(os.Stderr, err.Error())
			syscall.Exit(1)
		}
	}

	Core.New(&config)
}