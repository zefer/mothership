package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/fhs/gompd/mpd"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	port    = flag.String("port", ":8080", "listen port")
)

func main() {
	flag.Parse()
	glog.Infof("Starting API for MPD at %s.", *mpdAddr)

	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler)
	r.HandleFunc("/next", NextHandler)
	r.HandleFunc("/previous", PreviousHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist/")))
	http.Handle("/", r)
	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

func client() *mpd.Client {
	conn, err := mpd.Dial("tcp", *mpdAddr)
	if err != nil {
		glog.Errorln(err)
	}
	return conn
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	conn := client()
	defer conn.Close()
	data, err := conn.CurrentSong()
	if err != nil {
		glog.Errorln(err)
	}
	spew.Dump(data)
	b, err := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(b))
}

func NextHandler(w http.ResponseWriter, r *http.Request) {
	conn := client()
	defer conn.Close()
	err := conn.Next()
	if err != nil {
		glog.Errorln(err)
	}
	w.WriteHeader(http.StatusOK)
}

func PreviousHandler(w http.ResponseWriter, r *http.Request) {
	conn := client()
	defer conn.Close()
	err := conn.Previous()
	if err != nil {
		glog.Errorln(err)
	}
	w.WriteHeader(http.StatusOK)
}
