package main

import (
	"flag"

	"github.com/golang/glog"
)

var mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")

func main() {
	flag.Parse()
	glog.Infof("[main] Starting API for MPD at %s.", *mpdAddr)
}
