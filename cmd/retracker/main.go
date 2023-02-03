package main

import (
	"flag"
	"fmt"
	Core "github.com/zhukovra/retracker/core"
	"github.com/zhukovra/retracker/core/common"
	"os"
	"syscall"
)

const VERSION = 0.2

func PrintRepo() {
	fmt.Fprintln(os.Stderr, "\n# https://github.com/zhukovra/retracker")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		PrintRepo()
	}
	listen := flag.String("l", ":8080", "Listen address:port")
	age := flag.Float64("a", 180, "Keep 'n' minutes peer in memory")
	debug := flag.Bool("d", false, "Debug mode")
	xrealip := flag.Bool("x", false, "Get RemoteAddr from X-Real-IP header")
	ver := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *ver {
		fmt.Println(VERSION)
		PrintRepo()
		syscall.Exit(0)
	}

	config := common.Config{
		Listen:  *listen,
		Debug:   *debug,
		Age:     *age,
		XRealIP: *xrealip,
	}

	Core.New(&config)
}
