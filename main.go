package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/zefer/mpd-web/handlers"
	"github.com/zefer/mpd-web/mpd"
	"github.com/zefer/mpd-web/websocket"
)

var (
	client  *mpd.Client
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	port    = flag.String("port", ":8080", "listen port")
)

func sendStatus(c *websocket.Conn) {
	b, err := mpdStatus()
	if err != nil {
		return
	}
	c.Send(b)
}

func mpdStateChanged(subsystem string) {
	glog.Info("Changed subsystem:", subsystem)
	b, err := mpdStatus()
	if err != nil {
		glog.Errorln(err)
		return
	}
	websocket.Broadcast(b)
}

func main() {
	flag.Parse()
	glog.Infof("Starting API for MPD at %s.", *mpdAddr)

	// Send the browser the MPD state when they first connect.
	websocket.OnConnect(sendStatus)

	// This watcher notifies us when MPD's state changes, without polling.
	watch := mpd.NewWatcher(*mpdAddr)
	defer watch.Close()
	watch.OnStateChange(mpdStateChanged)

	// This client connection provides an API to MPD's commands.
	client = mpd.NewClient(*mpdAddr)
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/websocket", websocket.Serve)
	r.Handle("/next", handlers.NextHandler(client))
	r.Handle("/previous", handlers.PreviousHandler(client))
	r.Handle("/play", handlers.PlayHandler(client))
	r.Handle("/pause", handlers.PauseHandler(client))
	r.Handle("/randomOn", handlers.RandomOnHandler(client))
	r.Handle("/randomOff", handlers.RandomOffHandler(client))
	r.Handle("/files", handlers.FileListHandler(client))
	r.Handle("/playlist", handlers.PlayListHandler(client))
	r.Handle("/library/updated", handlers.LibraryUpdateHandler(client))

	// The front-end assets are served from a go-bindata file.
	r.PathPrefix("/").Handler(
		http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}),
	)
	http.Handle("/", r)
	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

func mpdStatus() ([]byte, error) {
	data, err := client.C.Status()
	if err != nil {
		return nil, err
	}
	song, err := client.C.CurrentSong()
	if err != nil {
		return nil, err
	}
	for k, v := range song {
		data[k] = v
	}
	b, err := json.Marshal(data)
	return b, err
}
